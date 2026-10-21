<!--
	The Design page's token editor (S-0383): one namespace of tokens, such as color or text,
	changed in place and bound to the edit model ($lib/design/tokens). A colour has a picker and a
	text field per theme, a length a number and its unit, a type step its font size and line
	height, a font family its stack as text; shadows, aliases, and other types are shown read-only.
	Text that is not a value flai reads stays in its field, marked, and makes no edit, so that
	every edit is one flai reads. Each colour shows the failing role pairs it takes part in, in the
	theme they fail in. The page owns the edits and recomputes edited and pairs from them; this
	writes only edits, pruned so that a value set back to its saved one leaves no edit.
-->
<script lang="ts">
	import type { DesignContrastPair, DesignTheme, DesignToken } from '$lib/server/repo';
	import { pairsOf, toSRGB, unsafeVariable } from '$lib/design/preview';
	import {
		changed,
		colorFromHex,
		colorText,
		cssNumber,
		editKey,
		parseColorText,
		parseEditKey,
		pruneEdits,
		valueProblem,
		type Edits,
		type Theme
	} from '$lib/design/tokens';

	let {
		namespace,
		saved,
		edited,
		edits = $bindable(),
		pairs,
		onedit
	}: {
		/** The tokens to edit: those whose path is namespace or starts with namespace and a dot. */
		namespace: string;
		/** The themes as flai last saved them, light then dark. */
		saved: readonly DesignTheme[];
		/** The saved themes with the edits applied: the page's editedThemes. */
		edited: readonly DesignTheme[];
		/** The edits, keyed by editKey, which this writes and prunes against saved. */
		edits: Edits;
		/** Both themes' role pairs on the edited tokens. */
		pairs: readonly DesignContrastPair[];
		/** Called with the edits after each change this makes to them. */
		onedit?: (edits: Edits) => void;
	} = $props();

	const uid = $props.id();
	const THEMES: Theme[] = ['light', 'dark'];
	const UNITS = ['px', 'rem'];

	type Dim = { value: number; unit: string };
	type Kind = 'color' | 'dimension' | 'typography' | 'font' | 'readonly';
	/** A token's value in one theme, or in both, as a light edit, when one field serves both. */
	type Cell = { theme: Theme; label: string; token: DesignToken };
	type Row = { path: string; kind: Kind; byTheme: Partial<Record<Theme, DesignToken>> };
	/** Text typed into a field, and what is wrong with it; null when it made its edit. */
	type Draft = { text: string; problem: string | null };
	type Parsed = { value: unknown } | { problem: string };
	/** A text field: its draft's key, its label, the value it shows, and how its text is read. */
	type Field = {
		cell: Cell;
		key: string;
		label: string;
		shown: string;
		parse: (text: string) => Parsed;
		width: string;
		describedBy?: string;
	};

	let drafts = $state<Record<string, Draft>>({});

	function kindOf(t: DesignToken): Kind {
		if (t.ref) return 'readonly';
		if (t.type === 'color') return 'color';
		if (t.type === 'dimension') return 'dimension';
		if (t.type === 'typography') return 'typography';
		if (t.type === 'fontFamily') return 'font';
		return 'readonly';
	}

	const inNamespace = (path: string) => path === namespace || path.startsWith(namespace + '.');

	// Every token path of the namespace, light's order first, each with its token in each theme.
	const rows = $derived.by(() => {
		const all: Row[] = [];
		for (const theme of THEMES)
			for (const token of edited.find((t) => t.name === theme)?.tokens ?? []) {
				if (!inNamespace(token.path)) continue;
				let row = all.find((r) => r.path === token.path);
				if (!row) all.push((row = { path: token.path, kind: kindOf(token), byTheme: {} }));
				row.byTheme[theme] = token;
			}
		return all;
	});

	const changedPaths = $derived(new Set(changed(saved, edits).map((t) => t.path)));

	const same = (a: DesignToken, b: DesignToken) =>
		a.css === b.css && JSON.stringify(a.variables) === JSON.stringify(b.variables);

	/** One cell when the row is not a colour and both themes render alike, else one per theme. */
	function cells(row: Row): Cell[] {
		const { light: l, dark: d } = row.byTheme;
		if (row.kind !== 'color' && l && d && same(l, d))
			return [{ theme: 'light', label: 'light and dark', token: l }];
		return THEMES.flatMap((theme) => {
			const token = row.byTheme[theme];
			return token ? [{ theme, label: theme, token }] : [];
		});
	}

	/** The value, when it may be written into a style; undefined when it is empty or refused. */
	const safe = (value: string | undefined): string | undefined =>
		value && unsafeVariable('--v', value) === null ? value : undefined;

	const keyOf = (c: Cell, part: string) => `${editKey(c.theme, c.token.path)}#${part}`;
	const domId = (key: string) => `${uid}-${key.replace(/[^\w-]/g, '-')}`;
	const labelOf = (c: Cell, what = '') => [c.token.path, c.label, what].filter(Boolean).join(', ');
	const problemOf = (key: string) => drafts[key]?.problem ?? null;

	function write(theme: Theme, path: string, value: unknown) {
		edits = pruneEdits(saved, { ...edits, [editKey(theme, path)]: value });
		onedit?.(edits);
	}

	/** Keeps the text typed into a field, and makes the edit it gives when flai reads it. */
	function typed(f: Field, text: string) {
		const parsed = f.parse(text);
		const { type, path } = f.cell.token;
		const problem = 'problem' in parsed ? parsed.problem : valueProblem(type, parsed.value, path);
		drafts[f.key] = { text, problem };
		if (!problem && 'value' in parsed) write(f.cell.theme, path, parsed.value);
	}

	/** The ids of what describes a field: its mark, and the colour's failing pairs. */
	function describedBy(f: Field): string | undefined {
		const ids = [problemOf(f.key) && domId(f.key), f.describedBy].filter(Boolean);
		return ids.length ? ids.join(' ') : undefined;
	}

	/** Lets a field that made its edit show its value again, as the edit model writes it. */
	function settle(key: string) {
		if (drafts[key] && !drafts[key].problem) delete drafts[key];
	}

	function picked(c: Cell, hex: string) {
		delete drafts[keyOf(c, 'colour')];
		const value = colorFromHex(hex, c.token.value);
		if (value) write(c.theme, c.token.path, value);
	}

	/** Drops every edit of the token, in both themes, and the text typed into its fields. */
	function reset(path: string) {
		const kept = Object.entries(edits).filter(([k]) => parseEditKey(k)?.path !== path);
		edits = pruneEdits(saved, Object.fromEntries(kept));
		for (const key of Object.keys(drafts))
			if (parseEditKey(key.slice(0, key.lastIndexOf('#')))?.path === path) delete drafts[key];
		onedit?.(edits);
	}

	const NUMBER = /^[+-]?(?:\d+\.?\d*|\.\d+)(?:e[+-]?\d+)?$/i;

	/** The value make gives of the number text holds. */
	function number(text: string, make: (n: number) => unknown): Parsed {
		const t = text.trim();
		if (!NUMBER.test(t)) return { problem: 'write a number, such as 1.25' };
		return { value: make(Number(t)) };
	}

	const NOT_A_COLOUR =
		'write a colour as #rrggbb, such as #1a2b3c, or as oklch(L C H), such as oklch(52% 0.18 255)';

	const QUOTES = ['"', "'"];

	/** A font family name, trimmed, without the quotes around it. */
	function unquote(name: string): string {
		const t = name.trim();
		const quoted = t.length > 1 && QUOTES.includes(t[0]) && t.at(-1) === t[0];
		return (quoted ? t.slice(1, -1) : t).trim();
	}

	/** A comma-separated font stack as its names, each unquoted, blanks dropped. */
	function fontStack(text: string): Parsed {
		const names = text.split(',').map(unquote);
		return { value: names.filter((n) => n !== '') };
	}

	const asDim = (v: unknown): Dim | null => {
		const d = v as Partial<Dim> | null;
		return d && typeof d.value === 'number' && typeof d.unit === 'string' ? (d as Dim) : null;
	};
	const fields = (v: unknown) => (v ?? {}) as Record<string, unknown>;

	function colourField(c: Cell): Field {
		const value = c.token.value;
		const parse = (t: string): Parsed => {
			const colour = parseColorText(t, value);
			return colour ? { value: colour } : { problem: NOT_A_COLOUR };
		};
		const key = keyOf(c, 'colour');
		const describedBy = failures(c).length ? domId(key + '-contrast') : undefined;
		const shown = colorText(value);
		return { cell: c, key, label: labelOf(c), shown, parse, width: 'w-44', describedBy };
	}

	function fontField(c: Cell): Field {
		const v = c.token.value;
		const shown = Array.isArray(v) ? v.join(', ') : typeof v === 'string' ? v : '';
		const key = keyOf(c, 'font');
		return { cell: c, key, label: labelOf(c), shown, parse: fontStack, width: 'w-96 max-w-full' };
	}

	function numberField(
		c: Cell,
		part: string,
		what: string,
		n: number,
		make: (n: number) => unknown
	): Field {
		const parse = (t: string) => number(t, make);
		const key = keyOf(c, part);
		return { cell: c, key, label: labelOf(c, what), shown: cssNumber(n), parse, width: 'w-20' };
	}

	/** A length's number. */
	const dimNumber = (c: Cell, part: string, what: string, dim: Dim, make: (d: Dim) => unknown) =>
		numberField(c, part, what, dim.value, (n) => make({ ...dim, value: n }));

	function unitChosen(c: Cell, dim: Dim, make: (d: Dim) => unknown, unit: string) {
		write(c.theme, c.token.path, make({ ...dim, unit }));
	}

	/** A type step's line height as a number. */
	const lineField = (c: Cell, v: Record<string, unknown>, lh: number) =>
		numberField(c, 'line', 'line height', lh, (n) => ({ ...v, lineHeight: n }));

	/** The failing role pairs the token takes part in, in the cell's theme, in words. */
	function failures(c: Cell): string[] {
		return pairsOf(c.token.path, pairs)
			.filter((p) => p.theme === c.theme)
			.map((p) => {
				const fg = p.foreground === c.token.path;
				const where = fg ? `on ${p.background}` : `under ${p.foreground}`;
				return `Fails ${p.threshold}:1 ${where} in ${p.theme}: ${p.ratio.toFixed(2)}:1`;
			});
	}

	const inputClass = 'rounded border border-line-strong bg-surface px-2 py-1 font-mono text-xs';
	const selectClass = 'rounded border border-line-strong bg-surface px-1 py-1 text-xs';
	const pickerClass = 'h-7 w-9 shrink-0 cursor-pointer rounded border border-line-strong';
	const swatchClass = 'inline-block size-6 shrink-0 rounded border border-line-strong';
	const alertClass = 'mt-1 rounded bg-danger-soft px-1.5 text-xs text-danger';
	const changedClass = 'rounded bg-surface px-1.5 text-xs text-primary';
</script>

{#snippet problem(key: string)}
	{#if problemOf(key)}
		<p id={domId(key)} class={alertClass} role="alert">{problemOf(key)}</p>
	{/if}
{/snippet}

{#snippet textField(f: Field)}
	{@const bad = problemOf(f.key)}
	<input
		type="text"
		class={[inputClass, f.width, bad && 'border-danger']}
		aria-label={f.label}
		aria-invalid={bad ? 'true' : undefined}
		aria-describedby={describedBy(f)}
		value={drafts[f.key]?.text ?? f.shown}
		oninput={(e) => typed(f, e.currentTarget.value)}
		onchange={() => settle(f.key)}
	/>
{/snippet}

{#snippet dimField(c: Cell, part: string, what: string, dim: Dim, make: (d: Dim) => unknown)}
	<div class="flex items-center gap-1">
		{@render textField(dimNumber(c, part, what, dim, make))}
		<select
			class={selectClass}
			aria-label={labelOf(c, what ? `${what} unit` : 'unit')}
			value={dim.unit}
			onchange={(e) => unitChosen(c, dim, make, e.currentTarget.value)}
		>
			{#each UNITS as unit (unit)}
				<option value={unit}>{unit}</option>
			{/each}
		</select>
	</div>
	{@render problem(keyOf(c, part))}
{/snippet}

{#snippet colourCell(c: Cell)}
	{@const fails = failures(c)}
	<div class="flex items-center gap-2">
		<span class={swatchClass} style:background-color={safe(c.token.css)} data-testid="swatch"
		></span>
		<input
			type="color"
			class={pickerClass}
			aria-label={labelOf(c, 'picker')}
			value={toSRGB(c.token.value)?.hex ?? '#000000'}
			oninput={(e) => picked(c, e.currentTarget.value)}
		/>
		{@render textField(colourField(c))}
	</div>
	{@render problem(keyOf(c, 'colour'))}
	{#if fails.length}
		<div id={domId(keyOf(c, 'colour') + '-contrast')}>
			{#each fails as fail (fail)}
				<p class="mt-1 text-xs text-danger" data-testid="contrast-fail">{fail}</p>
			{/each}
		</div>
	{/if}
{/snippet}

{#snippet typographyCell(c: Cell)}
	{@const v = fields(c.token.value)}
	{@const size = asDim(v.fontSize)}
	{@const lh = v.lineHeight}
	{@const lhDim = asDim(lh)}
	<div class="space-y-1">
		{#if size}
			{@render dimField(c, 'size', 'font size', size, (d) => ({ ...v, fontSize: d }))}
		{/if}
		{#if typeof lh === 'number'}
			{@render textField(lineField(c, v, lh))}
			{@render problem(keyOf(c, 'line'))}
		{:else if lhDim}
			{@render dimField(c, 'line', 'line height', lhDim, (d) => ({ ...v, lineHeight: d }))}
		{/if}
	</div>
{/snippet}

<div data-testid="token-editor" data-namespace={namespace}>
	{#each rows as row (row.path)}
		<div class="border-t border-line py-2 first:border-t-0" data-testid="edit-{row.path}">
			<div class="flex items-center gap-2">
				<span class="font-mono text-xs">{row.path}</span>
				{#if changedPaths.has(row.path)}
					<span class={changedClass} data-testid="changed">changed</span>
					<button
						type="button"
						class="text-xs underline"
						aria-label="Reset {row.path}"
						onclick={() => reset(row.path)}>Reset</button
					>
				{/if}
			</div>
			{#if row.kind === 'readonly'}
				{@const t = (row.byTheme.light ?? row.byTheme.dark)!}
				<p class="mt-1 text-xs text-muted" data-testid="readonly">
					{#if t.ref}
						Follows <code>{t.ref}</code>: change that token, or this one in the token file.
					{:else}
						A {t.type || 'untyped'} token is edited in the token file.
					{/if}
				</p>
				{#each cells(row) as c (c.theme)}
					<div class="text-xs">
						<span class="text-muted">{c.label}</span>
						<span class="font-mono break-all">{c.token.css}</span>
					</div>
				{/each}
			{:else}
				<div class="mt-1 flex flex-wrap gap-4">
					{#each cells(row) as c (c.theme)}
						<div class="min-w-0" data-cell-theme={c.label.replaceAll(' ', '-')}>
							<div class="mb-1 text-xs text-muted">{c.label}</div>
							{#if row.kind === 'color'}
								{@render colourCell(c)}
							{:else if row.kind === 'dimension'}
								{@const d = asDim(c.token.value)}
								{#if d}
									{@render dimField(c, 'value', '', d, (x) => x)}
								{/if}
							{:else if row.kind === 'typography'}
								{@render typographyCell(c)}
							{:else}
								{@render textField(fontField(c))}
								{@render problem(keyOf(c, 'font'))}
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/each}
</div>
