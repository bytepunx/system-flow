// Whether a flai on the host has this dashboard connected (ADR-0029), shared by the header badge
// and the banner every page shows when it has not: the project is read through that flai (S-0073),
// so without it there is nothing to show, and the page should say so and how to fix it. The budget
// banner reads it too, for the providers flai serve holds every agent start on (S-0534).
import { api } from '$lib/api';

/** Set by flai while it is older than the newest flai release in the project's history (S-0181). */
export type FlaiOutdated = {
	running: string;
	newest: string;
	upgrade: string;
	message: string;
};

/** A provider flai serve holds every agent start on, after a run ended on its usage limit (S-0534). */
export type BudgetHold = {
	provider: string;
	/** The limit the provider named, such as usage_limit_reached. */
	limit?: string;
	/** The HTTP status the provider refused with, such as 429. */
	status?: number;
	since: string;
	/** When the hold ends, as flai records a time. */
	until: string;
	/** Set while a probe runs to learn whether the limit has reset. */
	probing?: boolean;
};

export type HostFlaiStatus = {
	configured: boolean;
	connected: boolean;
	since?: string;
	/** The connected flai's version. */
	flai?: string;
	/** This dashboard's own version, whether or not a flai is connected (S-0544). */
	dashboard?: string;
	error?: string;
	/** The last flai this dashboard refused for the project as unsigned (S-0237); while none is
	 * connected, error says so and what to do, so it is not read as only "not connected". */
	refused?: { at: string; reason: string; flai: string };
	/** project.info, when the connected flai answered it. */
	info?: { flai_outdated?: FlaiOutdated };
	/** The providers flai serve holds every agent start on; absent while it holds none. */
	budget?: BudgetHold[];
};

class HostFlaiState {
	status = $state<HostFlaiStatus | null>(null);
	private timer: ReturnType<typeof setTimeout> | null = null;
	private users = 0;

	/** The project can be shown: a flai is connected and answers. */
	get usable(): boolean {
		return !!this.status?.connected && !this.status.error;
	}

	async refresh(): Promise<void> {
		try {
			const r = await api('/api/agent');
			if (r.ok) this.status = await r.json();
		} catch {
			// The dashboard itself cannot be reached; the browser says that better than we can.
		}
	}

	start(): void {
		if (this.users++ > 0) return;
		const tick = async () => {
			await this.refresh();
			// Look again sooner while it is away, so the page recovers soon after flai returns.
			this.timer = setTimeout(tick, this.usable ? 10000 : 3000);
		};
		void tick();
	}

	stop(): void {
		if (--this.users > 0) return;
		if (this.timer) clearTimeout(this.timer);
		this.timer = null;
	}
}

export const hostFlai = new HostFlaiState();
