// S-0663: DeploymentRow is one story namespace's row on the Deployments page: the story, its
// column, the namespace, the commit and whether it is the head, the last deploy and checks run
// with their outcomes, logs, and who asked, a badge per step, and the actions Deploy (or Set up
// deploys), Run checks, and Teardown once confirmed, each a callback the page runs, all disabled
// while busy. Each case of its examples file is rendered here, so the examples stay what it shows.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { readFileSync } from 'node:fs';
import type { KlusterDeployment } from '$lib/server/repo';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params?: Record<string, string>) =>
		params ? route.replace(/\[(?:\.\.\.)?(\w+)\]/g, (_, k: string) => params[k]) : route
}));

import DeploymentRow from './DeploymentRow.svelte';

type Props = { deployment: KlusterDeployment; deployable?: boolean; busy?: boolean };
type Example = { name: string; props: Props };

const FILE =
	'../../../../design/components/flaiover/src/lib/components/DeploymentRow.examples.json';
const raw = readFileSync(new URL(FILE, import.meta.url), 'utf8');
const EXAMPLES: Example[] = JSON.parse(raw).examples;
const example = (name: string) => {
	const e = EXAMPLES.find((x) => x.name === name);
	if (!e) throw new Error(`no example named ${name} in DeploymentRow.examples.json`);
	return structuredClone(e.props);
};
const CHECKED = 'In progress, deployed at its head and checked';
const BEHIND = 'In progress, behind its head';
const FAILED = 'Its deploy failed';
const REVIEW = 'In review, its checks failed';
const AGENT = 'Asked for by its agent, while a job runs and with no deploy set up';

let component: ReturnType<typeof mount> | undefined;

function render(props: Record<string, unknown>) {
	const list = document.createElement('ul');
	document.body.appendChild(list);
	component = mount(DeploymentRow, { target: list, props: props as Props });
	flushSync();
}
const q = (id: string) => document.querySelector<HTMLElement>(`[data-testid="${id}"]`);
const logLink = (within: ParentNode) =>
	within.querySelector<HTMLAnchorElement>('[data-testid="log-link"]')!;
const text = (id: string) => q(id)?.textContent?.replace(/\s+/g, ' ').trim();
const buttonNamed = (label: string, within: ParentNode = document) =>
	Array.from(within.querySelectorAll('button')).find((b) => b.textContent?.trim() === label);
const steps = () => Array.from(document.querySelectorAll<HTMLAnchorElement>('a[data-step]'));
const step = (name: string) => steps().find((s) => s.dataset.step === name)!;

afterEach(() => {
	if (component) unmount(component);
	component = undefined;
	document.body.innerHTML = '';
});

describe('DeploymentRow examples', () => {
	it('names every case the component is shown in', () => {
		expect(EXAMPLES.map((e) => e.name)).toEqual([CHECKED, BEHIND, FAILED, REVIEW, AGENT]);
	});

	it('shows a story in progress at its head and checked, with Run checks', () => {
		render(example(CHECKED));
		const link = q('deployment-who')!.querySelector('a')!;
		expect(link.textContent).toBe('S-0663');
		expect(link.getAttribute('href')).toBe('/items/S-0663');
		expect(text('deployment-who')).toContain('The Deployments page shows the namespace');
		expect(text('deployment-column')).toBe('in progress');
		expect(document.body.textContent).toContain('s-0663');
		expect(text('deployment-commit')).toBe('bb8cec4f4d2e');
		expect(text('deployment-at-head')).toBe('at its head');
		expect(q('deployment-behind')).toBeNull();
		expect(text('deployment-outcome')).toBe('passed');
		expect(text('deployment-asked-by')).toContain('asked by agent-S-0663');
		expect(text('deployment-checks-outcome')).toBe('passed');
		expect(text('deployment-checks-asked-by')).toBe('at bb8cec4f4d2e, asked by agent-S-0663');
		const log = '.flai-cache/serve/checks/S-0663-20261113T091700Z.log';
		expect(text('deployment-checks-log')).toContain(log);
		// S-0785: the checks run's log is a link to the log page, not a path to copy
		const checks = logLink(q('deployment-checks-log')!);
		expect(checks.textContent).toBe(log);
		expect(new URL(checks.href).pathname).toBe('/ops/log');
		expect(new URL(checks.href).searchParams.get('path')).toBe(log);
		// the last deploy's log is the folder of its steps' logs: Its log opens Logs, and no
		// folder is linked as a log
		expect(q('deployment-logs')!.getAttribute('title')).toBeNull();
		expect(q('deployment-failed-log')).toBeNull();
		expect(document.querySelectorAll('[data-testid="log-link"]')).toHaveLength(1);
		expect(q('deployment-checks-tail')!.textContent).toContain('check api passed');
		expect(buttonNamed('Deploy')).toBeUndefined();
		expect(buttonNamed('Run checks')).toBeDefined();
		expect(buttonNamed('Teardown')).toBeDefined();
		const badges = steps().map((s) => s.textContent);
		expect(badges).toEqual(['build · passed', 'load · passed', 'deploy · passed']);
	});

	it('shows a story behind its head with its head, Deploy, and no Run checks', () => {
		render(example(BEHIND));
		expect(text('deployment-commit')).toBe('7a4950a50c1d');
		expect(text('deployment-behind')).toBe('behind its head, ad3101cc7e8f');
		expect(q('deployment-at-head')).toBeNull();
		expect(buttonNamed('Deploy')).toBeDefined();
		expect(buttonNamed('Run checks')).toBeUndefined();
		expect(text('deployment-checks')).toBe('none since this deploy');
		expect(step('build-web').textContent).toBe('web · passed');
		expect(step('build-web').getAttribute('href')).toBe('#build');
	});

	it('shows a failed deploy with its step and error, and offers Deploy but not Run checks', () => {
		render(example(FAILED));
		expect(text('deployment-outcome')).toBe('failed, at its deploy step');
		expect(q('deployment-outcome')!.className).toContain('text-danger');
		expect(text('deployment-error')).toContain('the deploy step exited 1');
		// S-0785: the failed step's log, in its error and beside Its log, links the log page
		const log = '.flai-cache/kluster/logs/s-0586/deploy.log';
		expect(logLink(q('deployment-error')!).dataset.path).toBe(log);
		expect(text('deployment-failed-log')).toBe(`Its deploy step's log: ${log}`);
		expect(logLink(q('deployment-failed-log')!).dataset.path).toBe(log);
		expect(text('deployment-asked-by')).toContain('asked by a shell');
		expect(buttonNamed('Deploy')).toBeDefined();
		expect(buttonNamed('Run checks')).toBeUndefined();
		expect(step('deploy').dataset.outcome).toBe('failed');
		expect(step('deploy').className).toContain('bg-danger-soft');
	});

	it('shows a story in review whose checks failed, asked by flai serve', () => {
		render(example(REVIEW));
		expect(text('deployment-column')).toBe('in review');
		expect(text('deployment-asked-by')).toContain('asked by flai serve');
		expect(text('deployment-checks-outcome')).toBe('failed');
		expect(q('deployment-checks-outcome')!.className).toContain('text-danger');
		expect(text('deployment-checks-asked-by')).toContain('asked by flai serve');
		expect(q('deployment-checks-tail')!.textContent).toContain('GET /_ready answered 503');
		expect(buttonNamed('Run checks')).toBeDefined();
		expect(buttonNamed('Deploy')).toBeUndefined();
	});

	it('shows a running deploy its agent asked for, Set up deploys, and Teardown disabled', () => {
		render(example(AGENT));
		expect(text('deployment-asked-by')).toContain('asked by agent-S-0640');
		expect(text('deployment-job')).toContain('kluster deploy runs against it');
		expect(text('deployment-outcome')).toBe('running');
		expect(buttonNamed('Deploy')).toBeUndefined();
		expect(text('deployment-to-setup')).toBe('Set up deploys');
		expect(buttonNamed('Teardown')!.disabled).toBe(true);
		expect(step('build').dataset.outcome).toBe('running');
	});
});

describe('DeploymentRow actions', () => {
	it('deploys the story through onDeploy', () => {
		const onDeploy = vi.fn();
		render({ ...example(BEHIND), onDeploy });
		buttonNamed('Deploy')!.click();
		expect(onDeploy).toHaveBeenCalledWith('S-0581');
	});

	it('runs the checks through onChecks', () => {
		const onChecks = vi.fn();
		render({ ...example(CHECKED), onChecks });
		buttonNamed('Run checks')!.click();
		expect(onChecks).toHaveBeenCalledWith('S-0663');
	});

	it('tears down only once the operator confirms, and keeps it on Keep it', () => {
		const onTeardown = vi.fn();
		render({ ...example(CHECKED), onTeardown });
		expect(q('deployment-confirm-teardown')).toBeNull();
		buttonNamed('Teardown')!.click();
		flushSync();
		const confirm = q('deployment-confirm-teardown')!;
		expect(confirm.textContent).toContain("Tear down S-0663's namespace s-0663?");
		expect(buttonNamed('Teardown')!.disabled).toBe(true);
		buttonNamed('Keep it', confirm)!.click();
		flushSync();
		expect(q('deployment-confirm-teardown')).toBeNull();
		expect(onTeardown).not.toHaveBeenCalled();
		buttonNamed('Teardown')!.click();
		flushSync();
		buttonNamed('Tear it down', q('deployment-confirm-teardown')!)!.click();
		flushSync();
		expect(onTeardown).toHaveBeenCalledWith('S-0663');
		expect(q('deployment-confirm-teardown')).toBeNull();
	});

	it('opens Setup through onSetup when the block names no deploy', () => {
		const onSetup = vi.fn();
		const onDeploy = vi.fn();
		render({ ...example(FAILED), deployable: false, onSetup, onDeploy });
		expect(buttonNamed('Deploy')).toBeUndefined();
		const to = q('deployment-to-setup')!;
		expect(to.getAttribute('href')).toBe('#setup');
		to.click();
		expect(onSetup).toHaveBeenCalledOnce();
		expect(onDeploy).not.toHaveBeenCalled();
	});

	it('opens Logs at the failed step, Logs at a step, and Build at a build step', () => {
		const onLogs = vi.fn();
		const onBuild = vi.fn();
		render({ ...example(FAILED), onLogs, onBuild });
		const logs = q('deployment-logs')!;
		expect(logs.getAttribute('href')).toBe('#logs');
		logs.click();
		expect(onLogs).toHaveBeenLastCalledWith('s-0586', 'deploy');
		step('load').click();
		expect(onLogs).toHaveBeenLastCalledWith('s-0586', 'load');
		expect(step('load').getAttribute('href')).toBe('#logs');
		step('build').click();
		expect(onBuild).toHaveBeenCalledWith('s-0586', 'build');
		expect(onLogs).toHaveBeenCalledTimes(2);
	});

	it('opens Logs at the namespace when the deploy failed at no step', () => {
		const onLogs = vi.fn();
		render({ ...example(CHECKED), onLogs });
		q('deployment-logs')!.click();
		expect(onLogs).toHaveBeenCalledWith('s-0663', undefined);
	});

	it('disables Deploy, Run checks, and Teardown while busy', () => {
		const onDeploy = vi.fn();
		render({ ...example(BEHIND), busy: true, onDeploy });
		expect(buttonNamed('Deploy')!.disabled).toBe(true);
		expect(buttonNamed('Teardown')!.disabled).toBe(true);
		buttonNamed('Deploy')!.click();
		expect(onDeploy).not.toHaveBeenCalled();
		unmount(component!);
		component = undefined;
		document.body.innerHTML = '';

		render({ ...example(CHECKED), busy: true });
		expect(buttonNamed('Run checks')!.disabled).toBe(true);
		expect(buttonNamed('Teardown')!.disabled).toBe(true);
	});

	it('disables Tear it down once a job starts while the confirmation is open', () => {
		const onTeardown = vi.fn();
		const props = $state({ ...example(CHECKED), busy: false, onTeardown });
		render(props);
		buttonNamed('Teardown')!.click();
		flushSync();
		expect(buttonNamed('Tear it down')!.disabled).toBe(false);
		props.busy = true;
		flushSync();
		expect(buttonNamed('Tear it down')!.disabled).toBe(true);
		buttonNamed('Tear it down')!.click();
		expect(onTeardown).not.toHaveBeenCalled();
	});
});
