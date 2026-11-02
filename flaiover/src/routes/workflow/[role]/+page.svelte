<script lang="ts">
	// /workflow/<name> (S-0420): any role agent.status reports, the built-in ones and those the
	// project defines, on the role page, read from /api/roles/<name>; a role flai does not report says
	// so in its place. It is asked again when the role's activity document changes, when flai serve
	// says a run of the role started or ended, which changes no file, and when the path names
	// another role. A reload that fails says so above what was last shown.
	import { page } from '$app/state';
	import { untrack } from 'svelte';
	import { api } from '$lib/api';
	import { debounced, follow, listen } from '$lib/events';
	import RolePage from '$lib/components/RolePage.svelte';
	import { concernsRole, type RoleView } from '$lib/roles';

	const name = $derived(page.params.role ?? '');
	/** The role in the page's words: its name with its hyphens as spaces, its first letter capital. */
	const title = $derived(
		name ? `${name[0].toUpperCase()}${name.slice(1).replaceAll('-', ' ')}` : 'Role'
	);

	let view = $state<RoleView | null>(null);
	let error = $state<string | null>(null);
	let missing = $state(false);

	async function load(asked: string) {
		try {
			const r = await api(`/api/roles/${encodeURIComponent(asked)}`);
			const body = await r.json().catch(() => ({}));
			if (asked !== name) return;
			if (r.status === 404) {
				view = null;
				missing = true;
				error = body.error ?? `the project runs no role named ${asked}`;
				return;
			}
			if (!r.ok) throw new Error(body.error ?? r.statusText);
			view = body as RoleView;
			missing = false;
			error = null;
		} catch (e) {
			if (asked === name) error = e instanceof Error ? e.message : String(e);
		}
	}

	$effect(() => {
		const asked = name;
		return untrack(() => {
			view = null;
			missing = false;
			error = null;
			void load(asked);
			const again = debounced(() => void load(asked));
			const stops = [
				// of the narratives, only the role's own activity document is this page's
				follow(['narrative'], (changes) => {
					if (
						!view ||
						changes.some((c) => c.kind !== 'narrative' || c.path === view!.activity.path)
					)
						again();
				}),
				listen({
					agent: (id) => {
						if (!view || concernsRole(view.role, id)) again();
					}
				})
			];
			return () => {
				for (const stop of stops) stop();
				again.stop();
			};
		});
	});
</script>

<svelte:head><title>{title} · flaiover</title></svelte:head>

<h1 class="mb-1 text-2xl font-semibold">{title}</h1>
{#if view}
	<p class="mb-4 text-sm text-muted">
		Read from <code>{view.activity.path}</code> and from what flai serve knows of the {name}'s runs;
		its file is <code>{view.role.path}</code>.
	</p>
{/if}
{#if missing}
	<!-- flai's words name the roles the project does run -->
	<p class="text-sm text-muted" data-testid="role-not-found">{error}</p>
{:else}
	<!-- a reload that fails says so above what was last shown rather than in its place (S-0178) -->
	{#if error}
		<p
			class="mb-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			role="alert"
		>
			{error}
		</p>
	{/if}
	{#if view}
		<RolePage role={view.role} activity={view.activity} onchanged={() => void load(name)} />
	{:else if !error}
		<p class="text-sm text-muted">Loading…</p>
	{/if}
{/if}
