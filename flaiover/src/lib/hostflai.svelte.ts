// Whether a flai on the host has this dashboard connected (ADR-0029), shared by the header badge
// and the banner every page shows when it has not: the project is read through that flai (S-0073),
// so without it there is nothing to show, and the page should say so and how to fix it.
import { api } from '$lib/api';

/** Set by flai while it is older than the newest flai release in the project's history (S-0181). */
export type FlaiOutdated = {
	running: string;
	newest: string;
	upgrade: string;
	message: string;
};

export type HostFlaiStatus = {
	configured: boolean;
	connected: boolean;
	since?: string;
	flai?: string;
	error?: string;
	/** project.info, when the connected flai answered it. */
	info?: { flai_outdated?: FlaiOutdated };
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
