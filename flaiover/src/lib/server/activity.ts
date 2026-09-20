// Who is working on what, derived by flai on the host from the narratives (S-0042, S-0074).
// Nothing is written for presence: an agent that stops logging simply ages out.
import type { Repo } from './repo';

export type StreamActivity = {
	stream: string;
	title: string;
	agent: string;
	session: string;
	updated: string;
	age_seconds: number;
	status: string; // the story's state; "unknown" when the story is not found
	blocked: boolean; // the story or a task of it has an open blocked interval
	task?: { id: string; title: string }; // the task in progress
	last_log?: { at: string; text: string };
	path: string;
};

/** Age is counted here from `updated`, so an answer kept since the last change does not stand still. */
export async function activity(
	repo: Repo,
	now = new Date()
): Promise<{ streams: StreamActivity[] }> {
	const got = await repo.remember<{ streams: StreamActivity[] | null }>('activity', 'activity.get');
	return {
		streams: (got.streams ?? []).map((s) => ({
			...s,
			age_seconds: Math.max(0, Math.round((now.getTime() - Date.parse(s.updated)) / 1000)) || 0
		}))
	};
}
