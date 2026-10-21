<script lang="ts" module>
	import { editable, editLabel, fileRef } from '$lib/filelinks';

	/** A finding of the lint the fix left, as flai's refusal of the acceptance gives it (ADR-0140). */
	export type Finding = {
		tier: string;
		/** main, or story: the story's branch holds the file, and story names it. */
		checkout: string;
		story?: string;
		path: string;
		line?: number;
		column?: number;
		level?: string;
		rule?: string;
		alias?: string;
		message: string;
	};

	/** The findings of a refusal's data, or none when it has none. */
	export function findingsOf(data: unknown): Finding[] {
		const list = (data as { findings?: unknown } | undefined)?.findings;
		return Array.isArray(list) ? (list as Finding[]) : [];
	}

	// A main-checkout finding is saved without a story; a story's on its branch.
	const storyOf = (f: Finding) => (f.checkout === 'story' ? f.story : undefined);

	// The key a finding's file is saved under: its checkout's story, if any, and its path.
	const fileKey = (f: Finding) => `${storyOf(f) ?? ''}:${f.path}`;

	// The editor opens a markdown file under the layout's folders that flai lets it edit.
	const opens = (f: Finding) => !!fileRef(f.path) && editable(f.path);

	/** The keys of the files the editor can open that are not yet among the saved. */
	export function unsavedFiles(findings: Finding[], saved: string[]): string[] {
		return [...new Set(findings.filter(opens).map(fileKey))].filter((k) => !saved.includes(k));
	}
</script>

<script lang="ts">
	// The findings of an acceptance refused for lint the fix left (S-0378): each file and line opens
	// the editor on the finding's own checkout, the story's branch or the main checkout, and is
	// marked saved once saved (S-0465).
	import { openEditor } from '$lib/editor.svelte';

	let {
		id,
		status,
		findings,
		saved = [],
		onsaved
	}: {
		id: string;
		status?: string;
		findings: Finding[];
		/** The fileKey of each file saved so far. */
		saved?: string[];
		/** Called with a file's fileKey once flai has saved it. */
		onsaved?: (key: string) => void;
	} = $props();

	function openFinding(f: Finding) {
		const key = fileKey(f);
		openEditor({ path: f.path, line: f.line, story: storyOf(f), onsaved: () => onsaved?.(key) });
	}
</script>

<div
	class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger"
	role="alert"
	data-testid="accept-findings"
>
	<p class="font-medium">Not accepted.{status ? ` ${id} is ${status}.` : ''}</p>
	<p>The fix left lint in this markdown: correct each file at its line and save it.</p>
	<p>Accept is on again once each file is saved.</p>
	<ul class="mt-1 space-y-1 font-mono text-xs">
		{#each findings as f, i (i)}
			{@const where = f.line ? `${f.path}:${f.line}` : f.path}
			<li data-testid="accept-finding">
				{#if opens(f)}
					<button
						type="button"
						class="cursor-pointer underline"
						aria-label={editLabel({ path: f.path, line: f.line })}
						onclick={() => openFinding(f)}>{where}</button
					>
				{:else}
					<span class="font-medium">{where}</span>
				{/if}
				{#if storyOf(f)}<span class="text-muted">on {storyOf(f)}'s branch</span>{/if}
				{#if f.rule}<span>{f.alias ? `${f.rule}/${f.alias}` : f.rule}</span>{/if}
				<span class="whitespace-pre-wrap">{f.message}</span>
				{#if saved.includes(fileKey(f))}
					<span class="text-good" data-testid="finding-saved">saved</span>
				{/if}
			</li>
		{/each}
	</ul>
</div>
