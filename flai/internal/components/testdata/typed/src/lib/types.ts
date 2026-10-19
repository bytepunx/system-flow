// Types the fixture's components import through $lib.
export type Status = 'open' | 'closed';

export interface Story {
	id: string;
	points: number;
	done: boolean;
	size?: 'S' | 'M' | 'L';
	owner: string | null;
	due: Date;
	tags: string[];
	pair: [string, number];
	meta: Record<string, string>;
	onopen: () => void;
	close(): void;
}

export interface Tree {
	name: string;
	children: Tree[];
	parent?: Tree;
}

export type Deep = { a: { b: { c: { d: { e: { f: { g: string } } } } } } };

export type Props = { story: Story; count: number };
