// Which of the catalogue's components a story's branch changes (S-0442), from the paths flai's
// stream.diff answers: the Components page marks them and lists them first, and the review page
// offers Open in Components for a branch that changes one.
import type { Component, ComponentsCatalogue } from '$lib/server/repo';

/** A component as the catalogue names it: its sub-project's path and its own, both in slash form. */
export type ComponentRef = Pick<Component, 'project' | 'path'>;

const trimmed = (p: string) => p.replace(/^\.\/+/, '').replace(/\/+$/, '');

/** The repository-relative path of a component's file; a sub-project at the root adds nothing. */
export function componentFile(c: ComponentRef): string {
	const project = trimmed(c.project);
	const path = trimmed(c.path);
	return project && project !== '.' ? `${project}/${path}` : path;
}

/**
 * The paths a stream.diff answer changes, repository-relative: each file's path and, for a rename,
 * its old path; [] for an answer without files.
 */
export function diffPaths(diff: unknown): string[] {
	const files = (diff as { files?: unknown } | null)?.files;
	if (!Array.isArray(files)) return [];
	return files.flatMap((f: { path?: unknown; old_path?: unknown }) =>
		[f?.path, f?.old_path].filter((p): p is string => typeof p === 'string' && p !== '')
	);
}

/** The files, as componentFile gives them, of the catalogue's components that the paths change. */
export function changedComponents(
	catalogue: ComponentsCatalogue,
	paths: readonly string[]
): Set<string> {
	const changed = new Set(paths.map(trimmed));
	const files = new Set<string>();
	for (const p of catalogue.projects)
		for (const c of p.components) {
			const file = componentFile(c);
			if (changed.has(file)) files.add(file);
		}
	return files;
}

/** Whether any of the paths is a component of the catalogue. */
export function changesAComponent(
	catalogue: ComponentsCatalogue,
	paths: readonly string[]
): boolean {
	return changedComponents(catalogue, paths).size > 0;
}

/** The components with those changed first, each part in the order it was given. */
export function changedFirst<T extends ComponentRef>(
	components: readonly T[],
	changed: ReadonlySet<string>
): T[] {
	const isChanged = (c: T) => changed.has(componentFile(c));
	return [...components.filter(isChanged), ...components.filter((c) => !isChanged(c))];
}
