// Reads the props of Svelte 5 components for flai (S-0437). flai embeds this
// script and runs it with node in a sub-project's folder: it reads the
// components' paths, relative to that folder, as a JSON array on standard
// input, parses each with the sub-project's own svelte/compiler, and prints
// one JSON object on standard output:
//
//   {"components": [{"path", "props": [{"name", "kind", "type", "default",
//     "required", "literals"}], "error"}]}
//
// or, when the compiler cannot be loaded, {"error": {"code", "message"}}.
// It reads each component's source; it never runs one.
import fs from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';

const objectTypes = new Set(['Record', 'Partial', 'Required', 'Readonly', 'Pick', 'Omit', 'Object']);

const compiler = loadCompiler();
if (compiler.error) {
	print({ error: compiler.error });
} else {
	const files = JSON.parse(fs.readFileSync(0, 'utf8'));
	print({ components: files.map((file) => component(compiler.parse, file)) });
}

function print(answer) {
	process.stdout.write(JSON.stringify(answer) + '\n');
}

// loadCompiler resolves svelte/compiler from the sub-project, the working
// directory, not from this script's location.
function loadCompiler() {
	const require = createRequire(path.join(process.cwd(), 'package.json'));
	let version = '';
	try {
		version = require('svelte/package.json').version;
	} catch {
		return { error: { code: 'svelte-missing', message: 'cannot resolve svelte from ' + process.cwd() } };
	}
	if (Number.parseInt(version, 10) < 5) {
		return { error: { code: 'svelte-old', message: version } };
	}
	try {
		return { parse: require('svelte/compiler').parse };
	} catch (e) {
		return { error: { code: 'svelte-missing', message: String(e && e.message) } };
	}
}

// component is one file's entry: its props, or the error that kept them.
function component(parse, file) {
	let source;
	try {
		source = fs.readFileSync(file, 'utf8');
	} catch (e) {
		return { path: file, props: [], error: 'cannot read it: ' + e.message };
	}
	let ast;
	try {
		ast = parse(source, { modern: true });
	} catch (e) {
		return { path: file, props: [], error: parseError(e) };
	}
	try {
		return { path: file, props: props(source, ast) };
	} catch (e) {
		return { path: file, props: [], error: 'cannot read its props: ' + e.message };
	}
}

function parseError(e) {
	const first = String(e && e.message).split('\n')[0];
	if (e && e.start && e.start.line) {
		return 'line ' + e.start.line + ', column ' + e.start.column + ': ' + first;
	}
	return first;
}

// props reads the instance script's $props() declaration: the names it
// destructures, with their defaults, and the members of its type.
function props(source, ast) {
	const text = (node) => source.slice(node.start, node.end);
	const types = localTypes(ast);
	const decl = propsDeclarator(ast.instance);
	if (!decl) return [];
	const call = decl.init;
	const annotation =
		(decl.id.typeAnnotation && decl.id.typeAnnotation.typeAnnotation) ||
		firstParam(call.typeArguments) ||
		firstParam(call.typeParameters);
	const members = annotation ? membersOf(annotation, types, new Set()) : null;
	const byName = new Map();
	for (const m of members || []) {
		const name = memberName(m);
		if (name !== null && !byName.has(name)) byName.set(name, m);
	}

	const out = [];
	const seen = new Set();
	if (decl.id.type === 'ObjectPattern') {
		for (const p of decl.id.properties) {
			if (p.type !== 'Property') continue;
			const name = keyName(p.key);
			if (name === null) continue;
			const def = p.value.type === 'AssignmentPattern' ? p.value.right : null;
			const member = byName.get(name);
			out.push(prop(name, member, def, types, text));
			seen.add(name);
		}
	}
	for (const [name, member] of byName) {
		if (!seen.has(name)) out.push(prop(name, member, null, types, text));
	}
	return out;
}

function firstParam(args) {
	return args && args.params && args.params.length > 0 ? args.params[0] : null;
}

// propsDeclarator is the top-level declarator whose value is $props().
function propsDeclarator(script) {
	if (!script) return null;
	for (const stmt of script.content.body) {
		if (stmt.type !== 'VariableDeclaration') continue;
		for (const d of stmt.declarations) {
			const init = d.init;
			if (init && init.type === 'CallExpression' && init.callee.type === 'Identifier' && init.callee.name === '$props') {
				return d;
			}
		}
	}
	return null;
}

// localTypes are the type aliases and interfaces both scripts declare, by
// name.
function localTypes(ast) {
	const types = new Map();
	for (const script of [ast.module, ast.instance]) {
		if (!script) continue;
		for (let stmt of script.content.body) {
			if (stmt.type === 'ExportNamedDeclaration' && stmt.declaration) stmt = stmt.declaration;
			if (stmt.type === 'TSTypeAliasDeclaration' || stmt.type === 'TSInterfaceDeclaration') {
				types.set(stmt.id.name, stmt);
			}
		}
	}
	return types;
}

// membersOf is the members of an object type flai can see whole: a type
// literal, a local alias or interface, or an intersection of them; null
// when it cannot, as for an imported type.
function membersOf(node, types, seen) {
	node = unwrap(node);
	switch (node.type) {
		case 'TSTypeLiteral':
			return node.members;
		case 'TSIntersectionType': {
			const all = [];
			for (const t of node.types) {
				const m = membersOf(t, types, seen);
				if (m === null) return null;
				all.push(...m);
			}
			return all;
		}
		case 'TSTypeReference': {
			const decl = localType(node, types);
			if (!decl || seen.has(decl)) return null;
			seen.add(decl);
			if (decl.type === 'TSTypeAliasDeclaration') return membersOf(decl.typeAnnotation, types, seen);
			const all = [];
			for (const h of decl.extends || []) {
				const m = membersOf({ type: 'TSTypeReference', typeName: h.expression }, types, seen);
				if (m === null) return null;
				all.push(...m);
			}
			return all.concat(decl.body.body);
		}
		default:
			return null;
	}
}

function localType(ref, types) {
	return ref.typeName && ref.typeName.type === 'Identifier' ? types.get(ref.typeName.name) : undefined;
}

function unwrap(node) {
	while (node && (node.type === 'TSParenthesizedType' || node.type === 'TSTypeAnnotation')) {
		node = node.typeAnnotation;
	}
	return node;
}

function keyName(key) {
	if (key.type === 'Identifier') return key.name;
	if (key.type === 'Literal' && typeof key.value === 'string') return key.value;
	return null;
}

function memberName(m) {
	if (m.type !== 'TSPropertySignature' && m.type !== 'TSMethodSignature') return null;
	if (m.computed) return null;
	return keyName(m.key);
}

// prop is one prop's entry from its type member, if any, and its default,
// if any.
function prop(name, member, def, types, text) {
	let kind = 'unknown';
	let type = '';
	let literals;
	let optional = false;
	if (member) {
		optional = Boolean(member.optional);
		if (member.type === 'TSMethodSignature') {
			kind = 'function';
			type = text({ start: member.key.end, end: member.end })
				.replace(/^\s*\?/, '')
				.replace(/[;,]\s*$/, '')
				.trim();
		} else if (member.typeAnnotation) {
			const t = unwrap(member.typeAnnotation);
			type = text(t);
			({ kind, literals } = classify(t, types, new Set()));
		}
	} else if (def) {
		kind = kindOfValue(def);
	}
	const out = { name, kind, type, default: def ? text(def) : '', required: !optional && !def };
	if (literals) out.literals = literals;
	return out;
}

// classify gives a type's kind and, for a union of literals, its literals,
// each as JSON.
function classify(node, types, seen) {
	node = unwrap(node);
	switch (node.type) {
		case 'TSStringKeyword':
			return { kind: 'string' };
		case 'TSNumberKeyword':
			return { kind: 'number' };
		case 'TSBooleanKeyword':
			return { kind: 'boolean' };
		case 'TSArrayType':
		case 'TSTupleType':
			return { kind: 'array' };
		case 'TSTypeLiteral':
		case 'TSMappedType':
		case 'TSObjectKeyword':
			return { kind: 'object' };
		case 'TSFunctionType':
		case 'TSConstructorType':
			return { kind: 'function' };
		case 'TSTypeOperator':
			return node.operator === 'readonly' ? classify(node.typeAnnotation, types, seen) : { kind: 'unknown' };
		case 'TSLiteralType': {
			const l = literal(node);
			return l === null ? { kind: 'unknown' } : { kind: 'union', literals: [l] };
		}
		case 'TSUnionType':
			return classifyUnion(node, types, seen);
		case 'TSIntersectionType':
			return node.types.every((t) => classify(t, types, seen).kind === 'object') ? { kind: 'object' } : { kind: 'unknown' };
		case 'TSTypeReference':
			return classifyReference(node, types, seen);
		default:
			return { kind: 'unknown' };
	}
}

function classifyReference(node, types, seen) {
	const name = node.typeName.type === 'Identifier' ? node.typeName.name : node.typeName.right && node.typeName.right.name;
	const decl = localType(node, types);
	if (decl) {
		if (seen.has(decl)) return { kind: 'unknown' };
		seen.add(decl);
		return decl.type === 'TSInterfaceDeclaration' ? { kind: 'object' } : classify(decl.typeAnnotation, types, seen);
	}
	if (name === 'Snippet') return { kind: 'snippet' };
	if (name === 'Array' || name === 'ReadonlyArray') return { kind: 'array' };
	if (name === 'Function') return { kind: 'function' };
	if (node.typeName.type === 'Identifier' && objectTypes.has(name)) return { kind: 'object' };
	return { kind: 'unknown' };
}

// classifyUnion strips undefined and null, then gives the kind every member
// left shares, or union when each is a string, number, or boolean literal.
function classifyUnion(node, types, seen) {
	const rest = node.types.map(unwrap).filter((t) => t.type !== 'TSUndefinedKeyword' && t.type !== 'TSNullKeyword');
	if (rest.length === 0) return { kind: 'unknown' };
	if (rest.length === 1) return classify(rest[0], types, seen);
	const literals = rest.map((t) => (t.type === 'TSLiteralType' ? literal(t) : null));
	if (literals.every((l) => l !== null)) return { kind: 'union', literals };
	const kinds = rest.map((t) => classify(t, types, new Set(seen)).kind);
	if (kinds.every((k) => k === kinds[0] && k !== 'union')) return { kind: kinds[0] };
	return { kind: 'unknown' };
}

// literal is a string, number, or boolean literal type's value as JSON, or
// null for any other.
function literal(node) {
	const l = node.literal;
	if (l && l.type === 'Literal' && ['string', 'number', 'boolean'].includes(typeof l.value)) {
		return JSON.stringify(l.value);
	}
	if (l && l.type === 'UnaryExpression' && l.operator === '-' && l.argument.type === 'Literal' && typeof l.argument.value === 'number') {
		return JSON.stringify(-l.argument.value);
	}
	return null;
}

// kindOfValue is the kind of an untyped prop's default.
function kindOfValue(node) {
	switch (node.type) {
		case 'Literal':
			if (typeof node.value === 'string') return 'string';
			if (typeof node.value === 'number') return 'number';
			if (typeof node.value === 'boolean') return 'boolean';
			return 'unknown';
		case 'TemplateLiteral':
			return 'string';
		case 'ArrayExpression':
			return 'array';
		case 'ObjectExpression':
			return 'object';
		case 'ArrowFunctionExpression':
		case 'FunctionExpression':
			return 'function';
		default:
			return 'unknown';
	}
}
