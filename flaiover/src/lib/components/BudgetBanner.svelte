<script lang="ts">
	// A provider out of budget holds every agent start on it (S-0534, ADR-0209 §5): flai serve on the
	// host starts none until the limit resets, so every page names each provider held and until when.
	// Not dismissible, since it shows the host's state and goes by itself once the hold clears.
	import { hostFlai, type BudgetHold } from '$lib/hostflai.svelte';
	import { localTime } from '$lib/localtime';

	const held = $derived(hostFlai.status?.budget ?? []);

	function outOfBudget(hold: BudgetHold): string {
		const limit = hold.limit ? ` (${hold.limit})` : '';
		return `${hold.provider} is out of budget${limit} until ${localTime(hold.until)}`;
	}

	const HOLDS = '; flai serve holds every agent start on it until then.';
	const PROBING = ' A probe is running to learn whether the limit has reset.';

	function meanwhile(hold: BudgetHold): string {
		return hold.probing ? HOLDS + PROBING : HOLDS;
	}

	// S-0720: the hold can be cleared by hand on the host. A provider key may hold a space
	// (`claude-code meclod`), so it is single-quoted for the shell when it holds anything but the
	// characters a shell word needs no quoting for.
	function shellWord(word: string): string {
		return /^[A-Za-z0-9_@%+=:,./-]+$/.test(word) ? word : `'${word.replaceAll("'", `'\\''`)}'`;
	}

	const CLEAR = 'To clear the hold by hand, run ';
	const ON_HOST = ' on the host.';

	function clearCommand(hold: BudgetHold): string {
		return `flai serve budget clear ${shellWord(hold.provider)}`;
	}
</script>

{#if held.length}
	<div class="border-b border-warn bg-warn-soft px-4 py-3 text-sm" role="status" data-budget-banner>
		<div class="mx-auto max-w-6xl">
			{#each held as hold (hold.provider)}
				<p class="mt-1 first:mt-0" data-budget-provider={hold.provider}>
					<span class="font-medium">{outOfBudget(hold)}</span>{meanwhile(hold)}
					{CLEAR}<code>{clearCommand(hold)}</code>{ON_HOST}
				</p>
			{/each}
		</div>
	</div>
{/if}
