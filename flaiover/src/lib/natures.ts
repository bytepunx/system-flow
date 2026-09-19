// The five natures and what each means, as design/system/work-hierarchy.md defines them. A
// closed list: adding one takes an ADR. The "new" form shows the meanings (S-0059); the release
// consequence is part of what the operator is choosing (design/conventions/git.md, ADR-0025).
export const NATURES = [
	{
		name: 'feature',
		meaning: 'New capability that did not exist. Accepting it cuts a minor release.'
	},
	{
		name: 'improvement',
		meaning: 'Existing capability made better, faster, or clearer. A patch release.'
	},
	{
		name: 'remediation',
		meaning: 'Fixing something that is wrong, including defects and tech debt. A patch release.'
	},
	{
		name: 'research',
		meaning:
			'Producing knowledge: a spike or investigation with a written finding as the deliverable. Accepted and pushed, no release.'
	},
	{
		name: 'experiment',
		meaning:
			'Testing a hypothesis with a defined success measure; may be thrown away. Stays on its branch and is not accepted onto main.'
	}
] as const;

export const NATURE_NAMES: string[] = NATURES.map((n) => n.name);
