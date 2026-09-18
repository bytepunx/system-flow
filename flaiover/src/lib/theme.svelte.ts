// Theme selection (S-0044): 'system' follows prefers-color-scheme, 'light'
// and 'dark' are explicit. The choice is a per-browser convenience kept in
// localStorage; the tokens for both themes live in routes/layout.css.
export type Mode = 'system' | 'light' | 'dark';
const KEY = 'flaiover-theme';
const ORDER: Mode[] = ['system', 'light', 'dark'];

function stored(): Mode {
	try {
		const v = localStorage.getItem(KEY);
		return v === 'light' || v === 'dark' ? v : 'system';
	} catch {
		return 'system';
	}
}
function systemDark(): boolean {
	return typeof matchMedia === 'function' && matchMedia('(prefers-color-scheme: dark)').matches;
}

class ThemeState {
	mode = $state<Mode>('system');
	systemIsDark = $state(false);
	dark = $derived(this.mode === 'dark' || (this.mode === 'system' && this.systemIsDark));

	/** Read the stored choice, watch the system preference, and stamp <html>. */
	start(): () => void {
		this.mode = stored();
		this.systemIsDark = systemDark();
		const mq = typeof matchMedia === 'function' ? matchMedia('(prefers-color-scheme: dark)') : null;
		const onChange = (e: MediaQueryListEvent) => (this.systemIsDark = e.matches);
		mq?.addEventListener('change', onChange);
		this.apply();
		return () => mq?.removeEventListener('change', onChange);
	}
	apply(): void {
		if (typeof document === 'undefined') return;
		document.documentElement.setAttribute('data-theme', this.dark ? 'dark' : 'light');
	}
	set(mode: Mode): void {
		this.mode = mode;
		try {
			if (mode === 'system') localStorage.removeItem(KEY);
			else localStorage.setItem(KEY, mode);
		} catch {
			/* private windows: the choice lasts for the page */
		}
		this.apply();
	}
	cycle(): void {
		this.set(ORDER[(ORDER.indexOf(this.mode) + 1) % ORDER.length]);
	}
}

export const themeState = new ThemeState();
