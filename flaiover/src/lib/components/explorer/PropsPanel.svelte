<!--
	The picked component's props (S-0440): each with its type and its default as the source writes
	them, and the control its kind takes, which controls.ts chooses. Each control holds text, so that
	what the operator types stays until it parses; text that does not parse is shown with why, and a
	required prop without a value is marked. A function prop has no control: the harness logs its
	calls. The panel holds no texts of its own: each change, and the reset, is handed up whole.
-->
<script lang="ts">
	import type { ComponentProp } from '$lib/server/repo';
	import { controlFor, missing, parse, toTexts } from '$lib/explorer/controls';

	let {
		props,
		texts,
		onchange
	}: {
		/** The component's props, as components.list answers them. */
		props: readonly ComponentProp[];
		/** Each prop's control's text, by name; a prop without one has no value. */
		texts: Record<string, string>;
		/** Called with every prop's text after a change, or after the reset to the defaults. */
		onchange: (texts: Record<string, string>) => void;
	} = $props();

	const field = 'w-full rounded border border-line-strong bg-surface px-2 py-1';

	const rows = $derived(
		props.map((prop) => {
			const text = texts[prop.name];
			const parsed = parse(prop, text);
			return {
				prop,
				text,
				// as the source writes it, else the kind flai read
				type: prop.type ?? prop.kind,
				control: controlFor(prop),
				error: parsed.ok ? null : parsed.error,
				missing: missing(prop, text)
			};
		})
	);

	function set(name: string, text: string) {
		onchange({ ...texts, [name]: text });
	}

	const id = (name: string) => `props-panel-${name}`;
</script>

<section aria-labelledby="props-panel-heading" data-testid="props-panel">
	<div class="flex items-center justify-between gap-2">
		<h3 id="props-panel-heading" class="text-sm font-semibold">Props</h3>
		{#if props.length}
			<button
				type="button"
				class="rounded border border-line px-2 py-0.5 text-xs"
				data-testid="props-reset"
				onclick={() => onchange(toTexts([...props]))}>Reset</button
			>
		{/if}
	</div>
	{#if !props.length}
		<p class="mt-1 text-xs text-muted" data-testid="props-none">It takes no props.</p>
	{/if}
	<ul class="mt-1">
		{#each rows as { prop, text, type, control, error, missing: unset } (prop.name)}
			{@const described = error || unset ? `${id(prop.name)}-note` : undefined}
			<li class="border-b border-line py-2 last:border-b-0" data-testid="prop-{prop.name}">
				<div class="flex flex-wrap items-baseline gap-x-2 text-xs">
					<label for={id(prop.name)} class="font-mono text-sm font-medium">{prop.name}</label>
					{#if prop.required}
						<span class="text-danger" data-testid="prop-required">required</span>
					{/if}
					<span class="font-mono text-muted" data-testid="prop-type">{type}</span>
					{#if prop.default !== undefined}
						<span class="font-mono text-muted" data-testid="prop-default">= {prop.default}</span>
					{/if}
				</div>
				<div class="mt-1">
					{#if !control}
						<p id={id(prop.name)} class="text-xs text-muted" data-testid="prop-logged">
							Logged: each call is listed under Events.
						</p>
					{:else if control.kind === 'text'}
						<input
							id={id(prop.name)}
							type="text"
							class={field}
							value={text ?? ''}
							aria-invalid={!!error}
							aria-describedby={described}
							oninput={(e) => set(prop.name, e.currentTarget.value)}
						/>
					{:else if control.kind === 'number'}
						<input
							id={id(prop.name)}
							type="number"
							step="any"
							class={field}
							value={text ?? ''}
							aria-invalid={!!error}
							aria-describedby={described}
							oninput={(e) => set(prop.name, e.currentTarget.value)}
						/>
					{:else if control.kind === 'checkbox'}
						<div class="flex items-center gap-2 text-xs">
							<input
								id={id(prop.name)}
								type="checkbox"
								checked={text === 'true'}
								aria-describedby={described}
								onchange={(e) => set(prop.name, String(e.currentTarget.checked))}
							/>
							{#if text === undefined}
								<span class="text-muted">not set</span>
							{/if}
						</div>
					{:else if control.kind === 'select'}
						<select
							id={id(prop.name)}
							class={field}
							value={text ?? ''}
							aria-invalid={!!error}
							aria-describedby={described}
							onchange={(e) => set(prop.name, e.currentTarget.value)}
						>
							<option value="">(not set)</option>
							{#each control.options as o (o.text)}
								<option value={o.text}>{o.label}</option>
							{/each}
						</select>
					{:else}
						<textarea
							id={id(prop.name)}
							class={[field, control.kind === 'json' && 'font-mono text-xs']}
							rows={control.kind === 'json' ? 4 : 2}
							value={text ?? ''}
							placeholder={control.kind === 'json' ? 'JSON' : 'Text, rendered as the snippet'}
							aria-invalid={!!error}
							aria-describedby={described}
							oninput={(e) => set(prop.name, e.currentTarget.value)}></textarea>
					{/if}
				</div>
				{#if error}
					<p id={described} class="mt-1 text-xs text-danger" data-testid="prop-error">
						{error}; not sent.
					</p>
				{:else if unset}
					<p id={described} class="mt-1 text-xs text-danger" data-testid="prop-missing">
						Required: give it a value.
					</p>
				{/if}
			</li>
		{/each}
	</ul>
</section>
