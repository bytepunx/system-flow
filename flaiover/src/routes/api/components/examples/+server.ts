import {
	repo,
	RepoError,
	storyParam,
	type ComponentAnswer,
	type ComponentExample
} from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

type SaveBody = { project?: unknown; component?: unknown; example?: unknown; story?: unknown };

const isObject = (v: unknown): v is Record<string, unknown> =>
	typeof v === 'object' && v !== null && !Array.isArray(v);

/** The example the body gives, with ADR-0164's keys and ADR-0187's answers in its order and no
 * others, answers left out when there are none; 400 when it does not fit the format. The answers
 * need only be an array: flai check's components.examples holds each to the format on save. */
function exampleOf(v: unknown): ComponentExample {
	if (!isObject(v) || typeof v.name !== 'string' || !v.name.trim() || !isObject(v.props))
		throw new RepoError(400, 'example needs a name that is not empty and props, an object');
	if (v.theme !== undefined && v.theme !== 'light' && v.theme !== 'dark')
		throw new RepoError(400, `example's theme is light or dark, or left out, not ${v.theme}`);
	if (v.width !== undefined && !(Number.isInteger(v.width) && (v.width as number) >= 1))
		throw new RepoError(400, `example's width is a whole number of pixels, 1 or more, or left out`);
	if (v.answers !== undefined && !Array.isArray(v.answers))
		throw new RepoError(400, `example's answers are an array, or left out`);
	return {
		name: v.name.trim(),
		props: v.props,
		...(v.theme ? { theme: v.theme } : {}),
		...(v.width !== undefined ? { width: v.width as number } : {}),
		...(Array.isArray(v.answers) && v.answers.length
			? { answers: v.answers as ComponentAnswer[] }
			: {})
	};
}

/** The examples in an examples file's content, as the file writes them; 422 when it is not one. */
function examplesIn(path: string, content: string): unknown[] {
	const fix = "fix it as flai check's components.examples says, then save the example again";
	let parsed: unknown;
	try {
		parsed = JSON.parse(content);
	} catch (e) {
		throw new RepoError(
			422,
			`${path} is not JSON (${e instanceof Error ? e.message : String(e)}): ${fix}`
		);
	}
	if (!isObject(parsed) || !Array.isArray(parsed.examples))
		throw new RepoError(422, `${path} is not an object whose examples is an array: ${fix}`);
	return parsed.examples;
}

/** examples with example in place of the one of its name, else after them; the others are kept
 * as the file writes them, answers and all. */
function withExample(examples: unknown[], example: ComponentExample): unknown[] {
	const named = (e: unknown) => isObject(e) && e.name === example.name;
	return examples.some(named)
		? examples.map((e) => (named(e) ? example : e))
		: [...examples, example];
}

/** The repository path of the component's examples file, as components.list gives it, and whether
 * the catalogue gives it examples; 404 when the catalogue has no such component, 501 when the host
 * flai is too old to say where its examples go. */
async function examplesFile(project: string, component: string, story?: string) {
	const catalogue = await repo().components(story);
	const found = catalogue.projects
		.find((p) => p.name === project)
		?.components.find((c) => c.path === component);
	if (!found)
		throw new RepoError(
			404,
			`${project} has no component ${component}${story ? ` in ${story}'s worktree` : ''}: pick one the catalogue lists`
		);
	if (!found.examples_path)
		throw new RepoError(
			501,
			'the host flai is too old to keep examples: upgrade flai on the host with flai self-upgrade'
		);
	return { path: found.examples_path, held: (found.examples ?? []).length > 0 };
}

/**
 * PUT /api/components/examples { project, component, example: { name, props, theme?, width?,
 * answers? }, story? }: adds the example to the component's examples file, or replaces the one of
 * its name, and saves it through flai's file.save on the host (S-0443, ADR-0164; answers since
 * S-0476, ADR-0187): on the main checkout, or on the story's branch when story names one. project
 * is the sub-project's manifest name and component its path in the sub-project; the file is the one
 * components.list names. When the catalogue gives the component examples, the file and its hash are
 * read afresh through file.get; else the file is created with an empty hash, and when it exists
 * after all, holding none, flai's conflict gives its content and hash and the example is added to
 * that, once. The file is written as JSON with two-space indentation and a final newline, with
 * answers only when the example has any; the other examples are written as the file held them.
 * Answers file.save's { path, hash, committed, … } with the file's examples and the log. 400 for a
 * malformed body; 404 for a component the catalogue lacks; 409 { error, current, hash, diff } when
 * the file changed between the read and the save; 422 { error, findings } when flai refuses the
 * content, such as a components.examples finding, or { error } when the file is not an examples
 * file; 501 when the host flai is too old.
 */
export const PUT: RequestHandler = ({ request }) =>
	respond(async () => {
		const body = (await request.json().catch(() => ({}))) as SaveBody;
		const { project, component } = body;
		if (typeof project !== 'string' || !project || typeof component !== 'string' || !component)
			throw new RepoError(400, 'project, component, and example are required');
		const example = exampleOf(body.example);
		const story = storyParam(body.story);
		const { path, held } = await examplesFile(project, component, story);
		const message = `example "${example.name}" of ${component}`;
		const save = async (content: string, hash: string) => {
			const examples = withExample(examplesIn(path, content), example);
			const next = JSON.stringify({ examples }, null, 2) + '\n';
			const { data, warnings } = await repo().saveFile({
				path,
				content: next,
				hash,
				message,
				story
			});
			return { ...data, examples, log: warnings };
		};
		if (held) {
			const file = await repo().file(path, story);
			return save(file.content, file.hash);
		}
		try {
			return await save('{"examples":[]}', '');
		} catch (e) {
			const found = e instanceof RepoError && e.status === 409 ? e.data : undefined;
			if (typeof found?.current !== 'string' || typeof found.hash !== 'string') throw e;
			return save(found.current, found.hash);
		}
	});
