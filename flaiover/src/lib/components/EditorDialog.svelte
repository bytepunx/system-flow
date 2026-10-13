<script lang="ts">
	// The document editor in a dialog over the current page (S-0378): whatever openEditor names
	// opens here, at its line and on its story's branch, so a link in a report, a refusal, or a
	// thread is fixed without leaving the page. Mounted once, in the layout.
	import DocEditor from '$lib/components/DocEditor.svelte';
	import { closeEditor, editorState } from '$lib/editor.svelte';
	import type { Saved } from '$lib/edit';

	const target = $derived(editorState.target);

	// Unsaved text is only in this dialog: ask before closing it, as the /edit page asks before leaving.
	function close() {
		if (editorState.dirty && !confirm('Close without saving your changes?')) return;
		closeEditor();
	}

	function saved(s: Saved) {
		target?.onsaved?.(s);
	}
</script>

<svelte:window
	onkeydown={(e) => {
		if (target && e.key === 'Escape') close();
	}}
	onbeforeunload={(e) => {
		if (target && editorState.dirty) e.preventDefault();
	}}
/>

{#if target}
	<div
		class="fixed inset-0 z-50 flex items-start justify-center bg-ink/40 p-4"
		role="presentation"
		onclick={(e) => {
			if (e.target === e.currentTarget) close();
		}}
	>
		<div
			class="relative max-h-[calc(100vh-2rem)] w-full max-w-6xl overflow-y-auto rounded-lg border border-line bg-ground p-5 text-ink shadow-lg"
			role="dialog"
			aria-modal="true"
			aria-label="Edit {target.path}"
		>
			<button
				type="button"
				class="absolute top-3 right-3 rounded border border-line px-2 py-1 text-xs text-muted hover:text-ink"
				onclick={close}>Close</button
			>
			{#key `${target.path}|${target.story ?? ''}`}
				<DocEditor
					path={target.path}
					line={target.line}
					story={target.story}
					ondirty={(d) => (editorState.dirty = d)}
					onsaved={saved}
				/>
			{/key}
		</div>
	</div>
{/if}
