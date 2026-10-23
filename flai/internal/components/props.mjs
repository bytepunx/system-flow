// Reads the props of Svelte 5 components for flai (S-0437). flai embeds this
// script and runs it with node in a sub-project's folder: it reads the
// components' paths, relative to that folder, as a JSON array on standard
// input, parses each with the sub-project's own svelte/compiler, and prints
// one JSON object on standard output:
//
//   {"components": [{"path", "props": [{"name", "kind", "type", "default",
//     "required", "literals", "skeleton"}], "error"}]}
//
// or, when the compiler cannot be loaded, {"error": {"code", "message"}}.
// A prop of kind array, object, or unknown carries a skeleton, an empty value
// of its type, resolved with the sub-project's own typescript when it has one
// (S-0449). It reads each component's source; it never runs one.
//
// Its arguments are the cache's folder and the hash of this script.
// Resolving the skeletons costs about two seconds for a sub-project whose
// components import much; the cache keeps them, with a stamp of every path
// typescript looked at, so a run whose typings, typescript, and those paths
// are all as they were takes them from the cache without loading typescript
// (S-0478). Without a cache folder, every run resolves them.
import crypto from 'node:crypto';
import fs from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';

const objectTypes = new Set(['Record', 'Partial', 'Required', 'Readonly', 'Pick', 'Omit', 'Object']);

// fallbacks are the skeletons of the json kinds when a prop's type cannot be
// resolved; a prop of any other kind has no skeleton.
const fallbacks = { array: [], object: {}, unknown: null };

// propsAlias names the type typing gives the props; typedSkeletons reads it back.
const propsAlias = '__flai_props__';

// unresolved is the skeleton of a type that cannot be resolved, such as one
// imported from a module that does not resolve; omitted is a function's,
// whose key an object leaves out.
const unresolved = Symbol('unresolved');
const omitted = Symbol('omitted');

// maxDepth is how many objects and arrays deep a skeleton goes.
const maxDepth = 6;

// maxCached is how many runs' skeletons the cache keeps, one file each, so
// that runs over different components, such as the components list and the
// examples check, do not take each other's place.
const maxCached = 4;

const [cacheDir = '', scriptHash = ''] = process.argv.slice(1);

const require = createRequire(path.join(process.cwd(), 'package.json'));
const compiler = loadCompiler();
if (compiler.error) {
	print({ error: compiler.error });
} else {
	const files = JSON.parse(fs.readFileSync(0, 'utf8'));
	const read = files.map((file) => component(compiler.parse, file));
	addSkeletons(read);
	print({ components: read.map((r) => r.entry) });
}

function print(answer) {
	process.stdout.write(JSON.stringify(answer) + '\n');
}

// loadCompiler resolves svelte/compiler from the sub-project, the working
// directory, not from this script's location.
function loadCompiler() {
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

// component is one file's entry, its props or the error that kept them, and
// the TypeScript that types its props, if its props have a type.
function component(parse, file) {
	const failed = (error) => ({ entry: { path: file, props: [], error }, typing: null });
	let source;
	try {
		source = fs.readFileSync(file, 'utf8');
	} catch (e) {
		return failed('cannot read it: ' + e.message);
	}
	let ast;
	try {
		ast = parse(source, { modern: true });
	} catch (e) {
		return failed(parseError(e));
	}
	try {
		const { list, annotation } = props(source, ast);
		return { entry: { path: file, props: list }, typing: annotation === null ? null : typing(source, ast, annotation) };
	} catch (e) {
		return failed('cannot read its props: ' + e.message);
	}
}

// typing is a module of the component's scripts, module first, with the type
// of its props as propsAlias. export {} keeps each component's declarations
// to its own module.
function typing(source, ast, annotation) {
	const body = (script) => (script ? source.slice(script.content.start, script.content.end) : '');
	return [body(ast.module), body(ast.instance), 'type ' + propsAlias + ' = ' + annotation + ';', 'export {};', ''].join('\n');
}

function parseError(e) {
	const first = String(e && e.message).split('\n')[0];
	if (e && e.start && e.start.line) {
		return 'line ' + e.start.line + ', column ' + e.start.column + ': ' + first;
	}
	return first;
}

// props reads the instance script's $props() declaration: the names it
// destructures, with their defaults, and the members of its type, and its
// type as written, or null when it has none.
function props(source, ast) {
	const text = (node) => source.slice(node.start, node.end);
	const types = localTypes(ast);
	const decl = propsDeclarator(ast.instance);
	if (!decl) return { list: [], annotation: null };
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
	return { list: out, annotation: annotation ? text(unwrap(annotation)) : null };
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

// addSkeletons gives each prop of kind array, object, or unknown its
// skeleton: an empty value of its type, or its kind's fallback when its type
// cannot be resolved or the sub-project has no typescript.
function addSkeletons(read) {
	const typed = read.filter((r) => r.typing !== null && r.entry.props.some(isJSON));
	const resolved = typed.length > 0 ? cachedSkeletons(typed) : {};
	for (const r of read) {
		const skeletonOf = Object.hasOwn(resolved, r.entry.path) ? resolved[r.entry.path] : {};
		for (const p of r.entry.props) {
			if (!isJSON(p)) continue;
			p.skeleton = Object.hasOwn(skeletonOf, p.name) ? skeletonOf[p.name] : fallbacks[p.kind];
		}
	}
}

// isJSON is whether a prop is of a kind that carries a skeleton.
function isJSON(p) {
	return Object.hasOwn(fallbacks, p.kind);
}

// cachedSkeletons are the typed components' resolved skeletons, from the
// cache when its key and every stamp it holds are as they are now, otherwise
// resolved and written to the cache. The key is this script, the sub-project,
// the typescript it resolves, and each typed component's typing and json
// props; the stamps are those of the paths typescript looked at.
function cachedSkeletons(typed) {
	const tsPath = resolveTypescript();
	const key = crypto
		.createHash('sha256')
		.update(JSON.stringify([scriptHash, process.cwd(), tsPath, typed.map((r) => [r.entry.path, r.typing, r.entry.props.filter(isJSON).map((p) => p.name)])]))
		.digest('hex');
	const hit = readCache(key);
	if (hit) return hit;
	const ts = tsPath ? loadTypescript() : null;
	if (!ts) return {};
	const watched = watch(ts);
	const resolved = typedSkeletons(ts, typed);
	writeCache(key, { stamps: [[['file', tsPath], stamp('file', tsPath)], ...watched()], skeletons: resolved });
	return resolved;
}

// resolveTypescript is where the sub-project's typescript resolves, or empty
// when it has none.
function resolveTypescript() {
	try {
		return require.resolve('typescript');
	} catch {
		return '';
	}
}

// loadTypescript is the sub-project's own typescript, or null when it has
// none, or one without the compiler API this script uses.
function loadTypescript() {
	try {
		const ts = require('typescript');
		return typeof ts.createProgram === 'function' ? ts : null;
	} catch {
		return null;
	}
}

// readCache is the skeletons the cache holds for key, or null when it holds
// none, or a stamp has changed. A file it answers from is touched, so that
// writeCache keeps the files in use.
function readCache(key) {
	if (!cacheDir) return null;
	const file = path.join(cacheDir, key + '.json');
	let cache;
	try {
		cache = JSON.parse(fs.readFileSync(file, 'utf8'));
	} catch {
		return null;
	}
	if (!cache || !Array.isArray(cache.stamps) || typeof cache.skeletons !== 'object') return null;
	for (const [[kind, p], was] of cache.stamps) {
		if (stamp(kind, p) !== was) return null;
	}
	try {
		const now = new Date();
		fs.utimesSync(file, now, now);
	} catch {
		// Untouched, it is only the first to go when the folder is full.
	}
	return cache.skeletons;
}

// writeCache writes cache as its key's file, through a file of this
// process's own that it renames, so that a run beside it reads one or the
// other, and removes all but the newest maxCached of the folder's files. A
// cache that cannot be written is left as it was: the next run resolves the
// skeletons again, which is all a missing cache costs.
function writeCache(key, cache) {
	if (!cacheDir) return;
	const file = path.join(cacheDir, key + '.json');
	const tmp = file + '.' + process.pid + '.tmp';
	try {
		fs.mkdirSync(cacheDir, { recursive: true });
		fs.writeFileSync(tmp, JSON.stringify(cache));
		fs.renameSync(tmp, file);
		const kept = fs
			.readdirSync(cacheDir)
			.filter((n) => n.endsWith('.json'))
			.map((n) => ({ n, at: fs.statSync(path.join(cacheDir, n)).mtimeMs }))
			.sort((a, b) => b.at - a.at);
		for (const { n } of kept.slice(maxCached)) fs.rmSync(path.join(cacheDir, n), { force: true });
	} catch {
		fs.rmSync(tmp, { force: true });
	}
}

// watch records each path typescript's system looks at, from now on, with
// its stamp as it was before typescript looked, and gives the function that
// lists them, each as [[kind, path], stamp]. A file is stamped by its kind,
// time, and size, and a file that is missing, as most of module
// resolution's probes are, by its directory's listing instead, which a file
// made there changes; a directory asked to exist is stamped by its kind
// alone, and one listed by its time; a real path by what it resolves to.
function watch(ts) {
	const looked = new Set();
	const seen = new Map();
	const add = (kind, p, s) => {
		const k = kind + '\0' + p;
		if (!seen.has(k)) seen.set(k, [[kind, p], s === undefined ? stamp(kind, p) : s]);
	};
	const on = (method, kind) => {
		const f = ts.sys[method];
		if (typeof f !== 'function') return;
		ts.sys[method] = function (p, ...rest) {
			if (typeof p === 'string' && !looked.has(kind + '\0' + p)) {
				looked.add(kind + '\0' + p);
				const s = stamp(kind, p);
				if (kind === 'file' && s === '-') add('list', path.dirname(p));
				else add(kind, p, s);
			}
			return f.call(this, p, ...rest);
		};
	};
	on('fileExists', 'file');
	on('readFile', 'file');
	on('getFileSize', 'file');
	on('getModifiedTime', 'file');
	on('directoryExists', 'dir');
	on('getDirectories', 'list');
	on('realpath', 'real');
	return () => [...seen.values()];
}

// stamp is what a path is now, as kind sees it.
function stamp(kind, p) {
	if (kind === 'real') {
		try {
			return fs.realpathSync(p);
		} catch {
			return '-';
		}
	}
	const st = fs.statSync(p, { throwIfNoEntry: false });
	if (!st) return '-';
	const type = st.isDirectory() ? 'd' : 'f';
	if (kind === 'dir') return type;
	if (kind === 'list') return type + ':' + st.mtimeMs;
	return type + ':' + st.mtimeMs + ':' + st.size;
}

// typedSkeletons builds one program over the typed components, each as a
// module beside its file, and gives, by each one's path, its json props'
// skeletons by name, leaving out a prop whose type cannot be resolved.
function typedSkeletons(ts, typed) {
	const cwd = process.cwd();
	const name = (r) => posix(path.resolve(cwd, r.entry.path)) + '.__flai.ts';
	const modules = new Map(typed.map((r) => [name(r), r.typing]));
	const options = compilerOptions(ts, cwd);
	const program = ts.createProgram([...modules.keys()], options, compilerHost(ts, options, modules));
	const checker = program.getTypeChecker();
	const out = {};
	// An old typescript, whose checker has no public isArrayType, gives the fallbacks.
	if (typeof checker.isArrayType !== 'function') return out;
	for (const r of typed) {
		const file = program.getSourceFile(name(r));
		const alias = file && file.statements.find((s) => ts.isTypeAliasDeclaration(s) && s.name.text === propsAlias);
		if (!alias) continue;
		const propsType = checker.getTypeFromTypeNode(alias.type);
		const sk = skeletons(ts, checker, alias);
		const byName = {};
		for (const p of r.entry.props) {
			if (!isJSON(p)) continue;
			const sym = checker.getPropertyOfType(propsType, p.name);
			if (!sym) continue;
			const s = sk(checker.getTypeOfSymbolAtLocation(sym, alias), []);
			if (s !== unresolved && s !== omitted) byName[p.name] = s;
		}
		out[r.entry.path] = byName;
	}
	return out;
}

function posix(p) {
	return p.split(path.sep).join('/');
}

// compilerOptions are the sub-project's tsconfig.json's, when it has one,
// for emitting nothing. A SvelteKit sub-project whose options do not map
// $lib, as when .svelte-kit/tsconfig.json, which it extends, is not generated
// yet, maps it to src/lib, as SvelteKit does.
function compilerOptions(ts, cwd) {
	let options = {};
	const config = path.join(cwd, 'tsconfig.json');
	if (ts.sys.fileExists(config)) {
		const read = ts.readConfigFile(config, ts.sys.readFile);
		if (!read.error) options = ts.parseJsonConfigFileContent(read.config, ts.sys, cwd, undefined, config).options;
	}
	if (!(options.paths && options.paths.$lib) && isSvelteKit()) {
		const lib = posix(path.join(cwd, 'src', 'lib'));
		options.paths = { ...options.paths, $lib: [lib], '$lib/*': [lib + '/*'] };
	}
	return { ...options, noEmit: true, skipLibCheck: true };
}

function isSvelteKit() {
	try {
		const pkg = require('./package.json');
		return Boolean((pkg.dependencies || {})['@sveltejs/kit'] || (pkg.devDependencies || {})['@sveltejs/kit']);
	} catch {
		return false;
	}
}

// compilerHost reads the disk, but for the components' modules, which it
// serves from memory. Like tsc, it reads JSDoc only in JavaScript files,
// where it gives types.
function compilerHost(ts, options, modules) {
	const host = ts.createCompilerHost(options, true);
	if (ts.JSDocParsingMode) host.jsDocParsingMode = ts.JSDocParsingMode.ParseForTypeErrors;
	const { getSourceFile, fileExists, readFile } = host;
	host.getSourceFile = (file, language, ...rest) =>
		modules.has(file) ? ts.createSourceFile(file, modules.get(file), language, true) : getSourceFile.call(host, file, language, ...rest);
	host.fileExists = (file) => modules.has(file) || fileExists.call(host, file);
	host.readFile = (file) => (modules.has(file) ? modules.get(file) : readFile.call(host, file));
	return host;
}

// skeletons gives the function that builds a type's skeleton, given the
// objects and arrays it is inside.
function skeletons(ts, checker, at) {
	const F = ts.TypeFlags;
	const skeleton = (type, outer) => {
		const f = type.flags;
		if (f & (F.Any | F.Unknown | F.Never)) return unresolved;
		if (f & (F.Null | F.Undefined | F.Void)) return null;
		if (f & F.Boolean) return false;
		if (type.isUnion()) return union(type, outer);
		if (f & F.BooleanLiteral) return checker.typeToString(type) === 'true';
		if (type.isStringLiteral() || type.isNumberLiteral()) return type.value;
		if (f & F.StringLike) return '';
		if (f & (F.NumberLike | F.BigIntLike)) return 0;
		if (f & F.NonPrimitive) return {};
		if (f & (F.Object | F.Intersection)) return object(type, outer);
		return unresolved;
	};
	// union is the skeleton of its first member but undefined and null, false
	// for a boolean, or null when there is no other member.
	const union = (type, outer) => {
		const rest = type.types.filter((t) => !(t.flags & (F.Null | F.Undefined | F.Void)));
		if (rest.length === 0) return null;
		if (rest.every((t) => t.flags & F.BooleanLiteral)) return rest.length > 1 ? false : skeleton(rest[0], outer);
		return skeleton(rest[0], outer);
	};
	const item = (type, outer) => {
		const s = skeleton(type, outer);
		return s === unresolved || s === omitted ? null : s;
	};
	const object = (type, outer) => {
		if (type.symbol && type.symbol.name === 'Date') return '';
		const array = checker.isArrayType(type);
		const tuple = checker.isTupleType(type);
		if (outer.includes(type) || outer.length >= maxDepth) return array || tuple ? null : {};
		const inner = [...outer, type];
		if (array) {
			const s = item(checker.getTypeArguments(type)[0], inner);
			return s === null ? [] : [s];
		}
		if (tuple) {
			return checker
				.getTypeArguments(type)
				.slice(0, type.target.elementFlags.length)
				.map((t) => item(t, inner));
		}
		const members = checker.getPropertiesOfType(type).filter((p) => !/^__[@#]/.test(String(p.escapedName)));
		if (members.length === 0) {
			return checker.getSignaturesOfType(type, ts.SignatureKind.Call).length > 0 ? omitted : {};
		}
		const out = {};
		for (const m of members) {
			const s = skeleton(checker.getTypeOfSymbolAtLocation(m, at), inner);
			if (s === omitted) continue;
			out[m.name] = s === unresolved ? null : s;
		}
		return out;
	};
	return skeleton;
}
