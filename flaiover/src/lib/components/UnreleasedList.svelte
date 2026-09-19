<script lang="ts">
	// The components a no-release story touched: code that lands on main without a version
	// (ADR-0025). Shown under "No release" in the acceptance confirmation and on the review page.
	let { unreleased }: { unreleased?: { component: string; files: string[] }[] | null } = $props();
</script>

{#if unreleased?.length}
	<ul class="mt-1 ml-4 list-disc" data-testid="unreleased">
		{#each unreleased as u (u.component)}
			<li>
				<span class="font-medium">{u.component}</span> lands on main without a release
				<span class="text-muted">({u.files.length} file{u.files.length === 1 ? '' : 's'})</span>
			</li>
		{/each}
	</ul>
{/if}
