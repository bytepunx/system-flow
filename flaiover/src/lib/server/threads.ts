import { repo } from './repo';

/** The designer's name for entries posted from the dashboard: the manifest owner until identities arrive with the hub. */
export async function designer(): Promise<string> {
	const m = await repo().manifest();
	return (m as { owner?: string }).owner || 'designer';
}
