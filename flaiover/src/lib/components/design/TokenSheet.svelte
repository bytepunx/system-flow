<!--
	The design page's token sheet (S-0382): a project's tokens as flai's design.get answers them, laid
	out in the dashboard's own chrome. Colours are swatches with their contrast on the ground in each
	theme, then flai's role pairs; the type scale is sample text at each step; font stacks, spacing,
	radii, and shadows are specimens drawn to size. Swatches and specimens show the token's own
	values, written only through style directives and only when unsafeVariable lets them through.
	Light and dark stand side by side where their values differ.
-->
<script lang="ts">
	import type { DesignContrast, DesignTheme, DesignToken } from '$lib/server/repo';
	import { contrastOnGround, GROUND, toSRGB, unsafeVariable } from '$lib/design/preview';

	let {
		themes,
		contrast
	}: {
		/** The design system's themes, light then dark, as design.get answers them. */
		themes: readonly DesignTheme[];
		/** flai's contrast report of the colour role pairs, per theme. */
		contrast: readonly DesignContrast[];
	} = $props();

	type Theme = DesignTheme['name'];
	const THEMES: Theme[] = ['light', 'dark'];
	const LABEL: Record<Theme, string> = { light: 'Light', dark: 'Dark' };
	/** A token's value in one theme, or in both when they are the same. */
	type Cell = { label: string; theme: Theme; token: DesignToken };
	type Row = { path: string; byTheme: Partial<Record<Theme, DesignToken>> };

	const SPACING_STEPS = [1, 2, 4, 6, 8, 12, 16];
	const SAMPLE = 'The quick brown fox jumps over the lazy dog';

	/** The value, when it may be written into a style; undefined when it is empty or refused. */
	const safe = (value: string | undefined): string | undefined =>
		value && unsafeVariable('--v', value) === null ? value : undefined;

	const tokensOf = (name: Theme) => themes.find((t) => t.name === name)?.tokens ?? [];
	const light = $derived(tokensOf('light'));
	const dark = $derived(tokensOf('dark'));

	type Kind = 'color' | 'font' | 'type' | 'spacing' | 'radius' | 'shadow' | 'dimension' | 'other';
	function kindOf(t: DesignToken): Kind {
		if (t.type === 'color') return 'color';
		if (t.type === 'fontFamily') return 'font';
		if (t.type === 'typography' || t.path.startsWith('text.')) return 'type';
		if (t.path === 'spacing') return 'spacing';
		if (t.path.startsWith('radius.')) return 'radius';
		if (t.type === 'shadow') return 'shadow';
		if (t.type === 'dimension') return 'dimension';
		return 'other';
	}

	// Every token path, light's order first, each with its token in each theme.
	const rows = $derived.by(() => {
		const all: Row[] = [];
		for (const name of THEMES)
			for (const token of name === 'light' ? light : dark) {
				let row = all.find((r) => r.path === token.path);
				if (!row) all.push((row = { path: token.path, byTheme: {} }));
				row.byTheme[name] = token;
			}
		return (kind: Kind) => all.filter((r) => kindOf((r.byTheme.light ?? r.byTheme.dark)!) === kind);
	});

	const same = (a: DesignToken, b: DesignToken) =>
		a.css === b.css && JSON.stringify(a.variables) === JSON.stringify(b.variables);

	/** One cell when the themes agree, else one per theme that has the token. */
	function cells(row: Row): Cell[] {
		const { light: l, dark: d } = row.byTheme;
		if (l && d && same(l, d)) return [{ label: 'Light and dark', theme: 'light', token: l }];
		return THEMES.flatMap((theme) => {
			const token = row.byTheme[theme];
			return token ? [{ label: LABEL[theme], theme, token }] : [];
		});
	}

	const onGround = $derived({ light: contrastOnGround(light), dark: contrastOnGround(dark) });
	const hasGround = (theme: Theme) =>
		(theme === 'light' ? light : dark).some((t) => t.path === GROUND);

	/** Which WCAG 2 level a ratio reaches: 4.5 for text, 3 for large text and UI, or neither. */
	function grade(exact: number): { text: string; tone: string } {
		if (exact >= 4.5) return { text: '≥ 4.5 text', tone: 'text-good' };
		if (exact >= 3) return { text: '≥ 3 large text, UI', tone: 'text-warn' };
		return { text: 'under 3', tone: 'text-danger' };
	}

	/** The theme's value of a token path, for a specimen's backdrop. */
	const cssOf = (theme: Theme, path: string) =>
		(theme === 'light' ? light : dark).find((t) => t.path === path)?.css;

	const lineHeight = (t: DesignToken) =>
		t.variables.find((v) => v.name.endsWith('--line-height'))?.value;

	/** The spacing unit times n, written out when the unit is a plain length. */
	function times(t: DesignToken, n: number): string {
		const v = t.value as { value?: unknown; unit?: unknown } | null;
		if (v && typeof v.value === 'number' && typeof v.unit === 'string')
			return `${+(v.value * n).toFixed(4)}${v.unit}`;
		return `${n} × ${t.css}`;
	}

	// flai's role pairs, one row per pair, with each theme's result beside the other.
	const pairs = $derived.by(() => {
		const all: {
			foreground: string;
			background: string;
			byTheme: Partial<Record<Theme, DesignContrast['pairs'][number]>>;
		}[] = [];
		for (const c of contrast)
			for (const p of c.pairs) {
				let row = all.find((r) => r.foreground === p.foreground && r.background === p.background);
				if (!row)
					all.push((row = { foreground: p.foreground, background: p.background, byTheme: {} }));
				row.byTheme[c.theme] = p;
			}
		return all;
	});

	const section = 'mt-6 first:mt-0';
	const h2 = 'text-lg font-semibold';
	const note = 'mt-1 text-sm text-muted';
	const th = 'py-2 pr-4 text-left font-medium text-muted';
</script>

{#snippet themeHead(c: Cell)}
	<div class="mb-1 text-xs text-muted">{c.label}</div>
{/snippet}

<div data-testid="token-sheet">
	<section class={section} aria-labelledby="sheet-colours">
		<h2 id="sheet-colours" class={h2}>Colours</h2>
		<p class={note}>
			Each colour with its WCAG 2 contrast on <code>{GROUND}</code>, the ground, in each theme:
			4.5:1 is enough for text, 3:1 for large text and interface parts.
		</p>
		{#each THEMES.filter((t) => !hasGround(t)) as theme (theme)}
			<p class="mt-1 text-sm text-warn">
				{LABEL[theme]} has no <code>{GROUND}</code>: its contrast on the ground is not measured.
			</p>
		{/each}
		<table class="mt-3 w-full text-sm">
			<thead>
				<tr><th class={th}>Token</th><th class={th}>Light</th><th class={th}>Dark</th></tr>
			</thead>
			<tbody>
				{#each rows('color') as row (row.path)}
					<tr class="border-t border-line align-top" data-testid="swatch-{row.path}">
						<td class="py-2 pr-4 font-mono text-xs">{row.path}</td>
						{#each THEMES as theme (theme)}
							{@const token = row.byTheme[theme]}
							{@const c = onGround[theme].find((g) => g.path === row.path)}
							<td class="py-2 pr-4" data-sheet-theme={theme}>
								{#if token}
									{@const srgb = toSRGB(token.value)}
									{@const hex = srgb?.hex}
									<div class="flex items-start gap-2">
										<span
											class="mt-0.5 inline-block size-8 shrink-0 rounded border border-line-strong"
											style:background-color={safe(token.css)}
											data-testid="swatch"
										></span>
										<div class="min-w-0 text-xs">
											<div class="font-mono">{hex ?? token.css}</div>
											{#if hex && token.css.toLowerCase() !== hex}
												<div class="font-mono break-all text-muted">{token.css}</div>
											{/if}
											{#if c}
												{@const g = grade(c.exact)}
												<div data-testid="contrast">
													{#if row.path === GROUND}
														<span class="text-muted">the ground</span>
													{:else}
														{c.ratio.toFixed(2)}:1 <span class={g.tone}>{g.text}</span>
													{/if}
												</div>
											{/if}
											{#if srgb?.clipped}
												<div class="text-warn" data-testid="clipped">
													outside sRGB: measured as {srgb.hex}
												</div>
											{/if}
										</div>
									</div>
								{:else}
									<span class="text-xs text-muted">not in this theme</span>
								{/if}
							</td>
						{/each}
					</tr>
				{/each}
			</tbody>
		</table>
	</section>

	{#if pairs.length}
		<section class={section} aria-labelledby="sheet-pairs">
			<h2 id="sheet-pairs" class={h2}>Role pairs</h2>
			<p class={note}>flai's contrast of the roles that are used together, at its thresholds.</p>
			<table class="mt-3 w-full text-sm">
				<thead>
					<tr><th class={th}>Pair</th><th class={th}>Light</th><th class={th}>Dark</th></tr>
				</thead>
				<tbody>
					{#each pairs as row (`${row.foreground} ${row.background}`)}
						<tr
							class="border-t border-line"
							data-testid="pair-{row.foreground}-on-{row.background}"
						>
							<td class="py-2 pr-4 font-mono text-xs"
								>{row.foreground} <span class="text-muted">on</span> {row.background}</td
							>
							{#each THEMES as theme (theme)}
								{@const p = row.byTheme[theme]}
								<td class="py-2 pr-4 text-xs whitespace-nowrap" data-sheet-theme={theme}>
									{#if p}
										<span
											class="mr-1 inline-block rounded border border-line px-1.5 font-semibold"
											style:color={safe(p.foreground_hex)}
											style:background-color={safe(p.background_hex)}>Aa</span
										>
										{p.ratio.toFixed(2)}:1
										<span class={p.passes ? 'text-good' : 'text-danger'} data-testid="verdict"
											>{p.passes ? 'passes' : 'fails'} {p.threshold}</span
										>
										{#if p.clipped}<span class="text-warn">clipped</span>{/if}
									{:else}
										<span class="text-muted">not measured</span>
									{/if}
								</td>
							{/each}
						</tr>
					{/each}
				</tbody>
			</table>
		</section>
	{/if}

	{#if rows('type').length}
		<section class={section} aria-labelledby="sheet-type">
			<h2 id="sheet-type" class={h2}>Type scale</h2>
			<p class={note}>Each step at its size and line height, in the sans stack.</p>
			<div class="mt-3 space-y-3">
				{#each rows('type') as row (row.path)}
					<div class="border-t border-line pt-2" data-testid="type-{row.path}">
						<div class="font-mono text-xs text-muted">{row.path}</div>
						<div class="flex flex-wrap gap-4">
							{#each cells(row) as c (c.theme)}
								<div class="min-w-0 flex-1">
									{@render themeHead(c)}
									<p
										class="overflow-hidden text-ellipsis whitespace-nowrap"
										style:font-size={safe(c.token.css)}
										style:line-height={safe(lineHeight(c.token))}
										style:font-family={safe(cssOf(c.theme, 'font.sans'))}
										data-testid="type-sample"
									>
										{SAMPLE}
									</p>
									<div class="font-mono text-xs text-muted">
										{c.token.css}{#if lineHeight(c.token)}&nbsp;/ {lineHeight(c.token)}{/if}
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if rows('font').length}
		<section class={section} aria-labelledby="sheet-fonts">
			<h2 id="sheet-fonts" class={h2}>Font stacks</h2>
			<p class={note}>
				Shown as this browser resolves the stack: a family it lacks falls back to the next.
			</p>
			<div class="mt-3 space-y-3">
				{#each rows('font') as row (row.path)}
					<div class="border-t border-line pt-2" data-testid="font-{row.path}">
						<div class="font-mono text-xs text-muted">{row.path}</div>
						{#each cells(row) as c (c.theme)}
							{@render themeHead(c)}
							<p class="text-lg" style:font-family={safe(c.token.css)} data-testid="font-sample">
								{SAMPLE}
							</p>
							<div class="font-mono text-xs break-words text-muted">{c.token.css}</div>
						{/each}
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if rows('spacing').length}
		<section class={section} aria-labelledby="sheet-spacing">
			<h2 id="sheet-spacing" class={h2}>Spacing</h2>
			<p class={note}>The unit every spacing utility multiplies, drawn at a few multiples.</p>
			{#each rows('spacing') as row (row.path)}
				{#each cells(row) as c (c.theme)}
					<div class="mt-3" data-testid="spacing-{c.theme}">
						{@render themeHead(c)}
						{#each SPACING_STEPS as n (n)}
							<div class="flex items-center gap-3 py-0.5 text-xs">
								<span class="w-28 shrink-0 font-mono text-muted">{n} · {times(c.token, n)}</span>
								<span
									class="inline-block h-3 rounded-sm bg-primary"
									style:width={safe(`calc(${c.token.css} * ${n})`)}
									data-testid="spacing-bar"
								></span>
							</div>
						{/each}
					</div>
				{/each}
			{/each}
		</section>
	{/if}

	{#if rows('radius').length}
		<section class={section} aria-labelledby="sheet-radii">
			<h2 id="sheet-radii" class={h2}>Radii</h2>
			<div class="mt-3 flex flex-wrap gap-6">
				{#each rows('radius') as row (row.path)}
					<div data-testid="radius-{row.path}">
						<div class="font-mono text-xs text-muted">{row.path}</div>
						<div class="flex gap-3">
							{#each cells(row) as c (c.theme)}
								<div>
									{@render themeHead(c)}
									<div
										class="size-16 border-2 border-line-strong bg-raised"
										style:border-radius={safe(c.token.css)}
										data-testid="radius-specimen"
									></div>
									<div class="mt-1 font-mono text-xs text-muted">{c.token.css}</div>
								</div>
							{/each}
						</div>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if rows('shadow').length}
		<section class={section} aria-labelledby="sheet-shadows">
			<h2 id="sheet-shadows" class={h2}>Shadows</h2>
			<p class={note}>Each on its theme's surface over its ground.</p>
			<div class="mt-3 space-y-3">
				{#each rows('shadow') as row (row.path)}
					<div class="border-t border-line pt-2" data-testid="shadow-{row.path}">
						<div class="font-mono text-xs text-muted">{row.path}</div>
						<div class="flex flex-wrap gap-3">
							{#each THEMES as theme (theme)}
								{@const token = row.byTheme[theme]}
								{#if token}
									<div class="min-w-0 flex-1">
										<div class="mb-1 text-xs text-muted">{LABEL[theme]}</div>
										<div
											class="rounded border border-line p-6"
											style:background-color={safe(cssOf(theme, GROUND))}
										>
											<div
												class="h-12 rounded"
												style:background-color={safe(cssOf(theme, 'color.surface'))}
												style:box-shadow={safe(token.css)}
												data-testid="shadow-specimen"
											></div>
										</div>
										<div class="mt-1 font-mono text-xs break-words text-muted">{token.css}</div>
									</div>
								{/if}
							{/each}
						</div>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if rows('dimension').length}
		<section class={section} aria-labelledby="sheet-dimensions">
			<h2 id="sheet-dimensions" class={h2}>Other lengths</h2>
			{#each rows('dimension') as row (row.path)}
				<div class="mt-2" data-testid="dimension-{row.path}">
					<div class="font-mono text-xs text-muted">{row.path}</div>
					{#each cells(row) as c (c.theme)}
						<div class="flex items-center gap-3 py-0.5 text-xs">
							<span class="w-28 shrink-0 text-muted">{c.label} · {c.token.css}</span>
							<span class="inline-block h-3 rounded-sm bg-primary" style:width={safe(c.token.css)}
							></span>
						</div>
					{/each}
				</div>
			{/each}
		</section>
	{/if}

	{#if rows('other').length}
		<section class={section} aria-labelledby="sheet-other">
			<h2 id="sheet-other" class={h2}>Other tokens</h2>
			<table class="mt-3 w-full text-sm">
				<tbody>
					{#each rows('other') as row (row.path)}
						{#each cells(row) as c (c.theme)}
							<tr class="border-t border-line" data-testid="other-{row.path}">
								<td class="py-1 pr-4 font-mono text-xs">{row.path}</td>
								<td class="py-1 pr-4 text-xs text-muted">{c.label}</td>
								<td class="py-1 font-mono text-xs break-all">{c.token.css}</td>
							</tr>
						{/each}
					{/each}
				</tbody>
			</table>
		</section>
	{/if}
</div>
