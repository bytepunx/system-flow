// The document editor over any page (S-0378): a link in a report, a refusal, a thread, or a page
// opens a file at a line, on a story's branch when it names one, without leaving the page.
// EditorDialog, mounted once in the layout, shows whatever is open here.
import type { Saved } from './edit';

/** What to open: a repository path, the 1-based file line to put the caret on, and the story whose branch to save on. */
export type EditorTarget = {
	path: string;
	line?: number;
	story?: string;
	/** Called after flai accepted a save. */
	onsaved?: (saved: Saved) => void;
};

/** The document open in the editor dialog, or null when it is closed. */
export class EditorState {
	target = $state.raw<EditorTarget | null>(null);
}

export const editorState = new EditorState();

/** Open a document in the editor dialog over the current page. */
export function openEditor(target: EditorTarget): void {
	editorState.target = { ...target };
}

/** Close the editor dialog at once, unsaved changes or not; the dialog's own close asks first. */
export function closeEditor(): void {
	editorState.target = null;
}
