<!--
	The design preview's sample components (S-0382, ADR-0149 § 3): generic buttons, inputs, a form,
	a card, an alert per state, a table, navigation, and type specimens, styled only through the
	token vocabulary of roles and steps that gallery.css maps to the project's --ds-* values.
	PreviewFrame mounts one per theme inside its frame; nothing here uses the dashboard's classes or
	theme. Tailwind reads this file for gallery.css, comments included: a word here that names a
	utility adds it to that stylesheet.
-->
<script lang="ts">
	let { caption }: { caption: string } = $props();

	// Ids unique to this panel, since the frame holds one gallery per theme.
	const id = $props.id();
	const tabs = ['Overview', 'Activity', 'Settings'];
	let tab = $state(tabs[0]);

	const states = [
		{ name: 'Success', box: 'border-success bg-success-soft text-success', text: 'Saved.' },
		{
			name: 'Warning',
			box: 'border-warning bg-warning-soft text-warning',
			text: 'Your trial ends in three days.'
		},
		{
			name: 'Danger',
			box: 'border-danger bg-danger-soft text-danger',
			text: 'The payment was declined.'
		},
		{ name: 'Info', box: 'border-info bg-info-soft text-info', text: 'A new version is out.' }
	];

	const button =
		'rounded-md px-4 py-2 text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-focus disabled:cursor-not-allowed disabled:opacity-50';
	const field =
		'w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-text focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-focus';
	const label = 'mb-1 block text-sm font-medium text-text';
	const section = 'space-y-3';
	const heading = 'text-xs font-semibold tracking-wide text-text-muted uppercase';

	// The frame is a blank document: a sample link or form must not navigate it away.
	const stay = (e: Event) => e.preventDefault();
</script>

<div class="bg-background text-text space-y-8 p-6 font-sans text-base" data-gallery>
	<p class="text-text-muted text-sm font-semibold">{caption}</p>

	<section class={section} aria-labelledby="{id}-type">
		<h2 id="{id}-type" class={heading}>Type</h2>
		<h3 class="text-4xl font-bold">The quick brown fox</h3>
		<h4 class="text-2xl font-semibold">Jumps over the lazy dog</h4>
		<h5 class="text-xl font-semibold">A smaller heading</h5>
		<p class="max-w-prose">
			Body text sets the tone of a page. It should be comfortable to read at length, with a line
			height that keeps lines apart without letting a paragraph fall apart.
		</p>
		<p class="text-text-muted text-sm">
			Muted text for hints and metadata, and <code class="font-mono">monospace code</code>.
		</p>
		<p class="text-text-muted text-xs">The smallest step, for captions.</p>
	</section>

	<section class={section} aria-labelledby="{id}-buttons">
		<h2 id="{id}-buttons" class={heading}>Buttons</h2>
		<div class="flex flex-wrap gap-2">
			<button type="button" class="{button} bg-primary text-on-primary shadow-sm">Primary</button>
			<button type="button" class="{button} bg-secondary text-on-secondary">Secondary</button>
			<button type="button" class="{button} text-primary hover:bg-surface">Ghost</button>
			<button type="button" class="{button} border border-danger bg-danger-soft text-danger"
				>Danger</button
			>
			<button type="button" class="{button} bg-primary text-on-primary" disabled>Disabled</button>
		</div>
	</section>

	<section class={section} aria-labelledby="{id}-inputs">
		<h2 id="{id}-inputs" class={heading}>Inputs</h2>
		<div class="grid gap-4 sm:grid-cols-2">
			<div>
				<label class={label} for="{id}-name">Name</label>
				<input id="{id}-name" class={field} type="text" placeholder="Ada Lovelace" />
			</div>
			<div>
				<label class={label} for="{id}-plan">Plan</label>
				<select id="{id}-plan" class={field}>
					<option>Free</option>
					<option>Team</option>
					<option>Enterprise</option>
				</select>
			</div>
			<div>
				<label class={label} for="{id}-email">Email</label>
				<input
					id="{id}-email"
					class="{field} border-danger"
					type="email"
					value="ada@"
					aria-invalid="true"
					aria-describedby="{id}-email-error"
				/>
				<p id="{id}-email-error" class="mt-1 text-sm text-danger">Enter a whole email address.</p>
			</div>
			<fieldset class="space-y-2">
				<legend class={label}>Notify me</legend>
				<label class="flex items-center gap-2 text-sm">
					<input type="checkbox" class="accent-primary" checked /> By email
				</label>
				<label class="flex items-center gap-2 text-sm">
					<input type="radio" name="{id}-often" class="accent-primary" checked /> Daily
				</label>
				<label class="flex items-center gap-2 text-sm">
					<input type="radio" name="{id}-often" class="accent-primary" /> Weekly
				</label>
			</fieldset>
			<div class="sm:col-span-2">
				<label class={label} for="{id}-notes">Notes</label>
				<textarea id="{id}-notes" class={field} rows="3">A longer answer.</textarea>
			</div>
		</div>
	</section>

	<section class={section} aria-labelledby="{id}-form">
		<h2 id="{id}-form" class={heading}>Form and card</h2>
		<form
			class="border-border max-w-md space-y-4 rounded-lg border bg-surface p-6 shadow-md"
			onsubmit={stay}
		>
			<h3 class="text-lg font-semibold">Sign in</h3>
			<div>
				<label class={label} for="{id}-user">Username</label>
				<input id="{id}-user" class={field} type="text" autocomplete="off" />
			</div>
			<div>
				<label class={label} for="{id}-password">Password</label>
				<input id="{id}-password" class={field} type="password" autocomplete="off" />
			</div>
			<div class="flex items-center justify-between gap-2">
				<button type="submit" class="{button} bg-primary text-on-primary">Sign in</button>
				<button type="button" class="{button} text-primary">Forgot password?</button>
			</div>
		</form>
		<article class="border-border max-w-md rounded-lg border bg-surface p-4 shadow-lg">
			<h3 class="text-lg font-semibold">A card</h3>
			<p class="text-text-muted text-sm">Updated two hours ago</p>
			<p class="mt-2">A card holds one subject: a summary, its details, and what to do next.</p>
		</article>
	</section>

	<section class={section} aria-labelledby="{id}-alerts">
		<h2 id="{id}-alerts" class={heading}>Alerts</h2>
		{#each states as s (s.name)}
			<div class="rounded-md border p-3 text-sm {s.box}">
				<strong class="font-semibold">{s.name}.</strong>
				{s.text}
			</div>
		{/each}
	</section>

	<section class={section} aria-labelledby="{id}-table">
		<h2 id="{id}-table" class={heading}>Table</h2>
		<table class="w-full border-collapse text-left text-sm">
			<caption class="text-text-muted pb-2 text-left">Invoices this month</caption>
			<thead>
				<tr class="border-border border-b">
					<th scope="col" class="py-2 pr-4 font-semibold">Invoice</th>
					<th scope="col" class="py-2 pr-4 font-semibold">Customer</th>
					<th scope="col" class="py-2 font-semibold">Status</th>
				</tr>
			</thead>
			<tbody>
				<tr class="border-border border-b">
					<td class="py-2 pr-4 font-mono">INV-001</td>
					<td class="py-2 pr-4">Acme</td>
					<td class="text-success py-2">Paid</td>
				</tr>
				<tr class="border-border border-b bg-surface">
					<td class="py-2 pr-4 font-mono">INV-002</td>
					<td class="py-2 pr-4">Globex</td>
					<td class="text-warning py-2">Due</td>
				</tr>
				<tr>
					<td class="py-2 pr-4 font-mono">INV-003</td>
					<td class="py-2 pr-4">Initech</td>
					<td class="py-2 text-danger">Overdue</td>
				</tr>
			</tbody>
		</table>
	</section>

	<section class={section} aria-labelledby="{id}-nav">
		<h2 id="{id}-nav" class={heading}>Navigation</h2>
		<nav
			class="flex items-center gap-4 rounded-md bg-surface px-4 py-2 text-sm"
			aria-label="Sample site"
		>
			<span class="font-semibold">Product</span>
			<a href="#{id}-nav" class="font-medium text-primary" aria-current="page" onclick={stay}
				>Home</a
			>
			<a href="#{id}-nav" class="text-text-muted hover:text-text" onclick={stay}>Pricing</a>
			<a href="#{id}-nav" class="text-text-muted hover:text-text" onclick={stay}>Docs</a>
		</nav>
		<div role="tablist" aria-label="Sample tabs" class="border-border flex gap-1 border-b">
			{#each tabs as t (t)}
				<button
					type="button"
					role="tab"
					aria-selected={tab === t}
					class="-mb-px border-b-2 px-3 py-2 text-sm {tab === t
						? 'border-primary font-medium text-primary'
						: 'text-text-muted hover:text-text border-transparent'}"
					onclick={() => (tab = t)}>{t}</button
				>
			{/each}
		</div>
		<div role="tabpanel" tabindex="0" aria-label={tab} class="text-sm" data-tab-panel>{tab}</div>
		<nav aria-label="Sample breadcrumbs" class="text-sm">
			<ol class="text-text-muted flex flex-wrap items-center gap-2">
				<li><a href="#{id}-nav" class="text-primary" onclick={stay}>Home</a></li>
				<li aria-hidden="true">/</li>
				<li><a href="#{id}-nav" class="text-primary" onclick={stay}>Projects</a></li>
				<li aria-hidden="true">/</li>
				<li aria-current="page" class="text-text">Design</li>
			</ol>
		</nav>
	</section>
</div>
