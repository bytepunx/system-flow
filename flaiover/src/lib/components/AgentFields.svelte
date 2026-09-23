<script lang="ts">
	// A story's agent (S-0103): the harness that runs it, the model, and options for the harness, one
	// key=value a line. The project's default shows as placeholders; what is typed here overrides it
	// for this story.
	import { agentLine, type Agent } from '$lib/agent';

	let {
		harness = $bindable(''),
		model = $bindable(''),
		config = $bindable(''),
		defaults,
		note,
		isDefault = false
	}: {
		harness?: string;
		model?: string;
		config?: string;
		defaults?: Agent;
		note?: string;
		/** Editing the project's default itself (S-0105): no "project default" line under it. */
		isDefault?: boolean;
	} = $props();

	const defaultConfig = $derived(
		Object.keys(defaults?.config ?? {})
			.sort()
			.map((k) => `${k}=${defaults!.config![k]}`)
			.join('\n')
	);
</script>

<fieldset class="rounded border border-line p-3 text-sm" data-testid="agent-fields">
	<legend class="px-1 text-xs text-muted">agent: who works this story</legend>
	<div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
		<label class="block">
			<span class="mb-1 block text-xs text-muted">harness</span>
			<input
				class="w-full rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={harness}
				placeholder={defaults?.harness ?? 'claude-code'}
				data-testid="agent-harness"
			/>
		</label>
		<label class="block">
			<span class="mb-1 block text-xs text-muted">model</span>
			<input
				class="w-full rounded border border-line-strong bg-surface px-2 py-1"
				bind:value={model}
				placeholder={defaults?.model ?? 'claude-opus-5-5'}
				data-testid="agent-model"
			/>
		</label>
		<label class="block">
			<span class="mb-1 block text-xs text-muted">config, one key=value a line</span>
			<textarea
				class="h-16 w-full rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs"
				bind:value={config}
				placeholder={defaultConfig || 'effort=high'}
				data-testid="agent-config"></textarea>
		</label>
	</div>
	<p class="mt-2 text-xs text-muted" data-testid="agent-default">
		{isDefault ? (note ?? '') : `project default: ${agentLine(defaults)}${note ? `. ${note}` : ''}`}
	</p>
</fieldset>
