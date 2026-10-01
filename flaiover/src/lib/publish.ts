/**
 * How this clone's release tags stand against its remote's (S-0174, flai's release.RemoteTags):
 * set on the publish preview only when there is something to say. `behind` names the components
 * whose tags here lag the remote's, when the preview offers no plan; `unchecked` says why the
 * remote could not be asked, when the plan is this clone's tags alone.
 */
export type RemoteTags = {
	remote: string;
	behind?: { component: string; local?: string; remote: string }[];
	unchecked?: string;
	fix: string;
	message: string;
};

/** An accepted item no plan covers, with why (I-0024, flai's release.Unplanned). */
export type Unplanned = { id: string; title?: string; reason: string };

/** Whether the remote has release tags newer than this clone's. */
export const lagging = (remote: RemoteTags | null | undefined): boolean =>
	(remote?.behind?.length ?? 0) > 0;

/**
 * The done lane's cards. flai keeps an archived card there only while its release looks pending;
 * while the clone lags its remote's tags that cannot be told from published, so those cards are
 * left out rather than listed as waiting (S-0174).
 */
export function doneLane<T extends { archived?: boolean }>(
	cards: T[],
	remote: RemoteTags | null | undefined
): T[] {
	return lagging(remote) ? cards.filter((c) => !c.archived) : cards;
}
