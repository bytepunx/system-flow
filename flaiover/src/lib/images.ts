// The project's agent images and flai serve's image jobs as flai's images methods answer them
// (S-0757, ADR-0295 §5), for the routes under /api/images and the pages and components that show
// them.

/** What an agent image is to the flai that must run in it (Go's agentimage.State). */
export type AgentImageState = 'missing' | 'current' | 'stale' | 'own';

/** One reason an image is stale, in a sentence (Go's agentimage.Reason). */
export type AgentImageReason = { kind: 'flai' | 'recipe' | 'base'; message: string };

/** The labels that tie an image to what it was built from (Go's agentimage.Labels). */
export type AgentImageLabels = { flai?: string; recipe?: string; base?: string; commit?: string };

/** What an image's labels must hold to be current; absent for a label not compared. */
export type AgentImageWant = { flai?: string; recipe?: string; base?: string };

/**
 * One image the project's agents need, as flai image status judges it (Go's agentimage.Status): a
 * harness image or the project's layer over one, for an inner harness, and whether it is on the
 * host, its labels there, its state, and why it is stale.
 */
export type AgentImage = {
	kind: 'harness' | 'layer';
	/** The inner harness it is for, such as claude-code. */
	inner: string;
	ref: string;
	/** Where its name came from: the manifest's layer, the host's own, or flai's default. */
	source: 'layer' | 'host' | 'default';
	/** A layer's Dockerfile, from the repository root. */
	dockerfile?: string;
	/** A layer's build arguments, from agent.image.args. */
	args?: Record<string, string>;
	present: boolean;
	/** Its ID on the host, when present. */
	id?: string;
	labels: AgentImageLabels;
	state: AgentImageState;
	reasons?: AgentImageReason[];
	want: AgentImageWant;
};

/** How an image job left one image (Go's agentimage.Way). */
export type ImageWayName = 'pulled' | 'built' | 'layered' | 'current' | 'left' | 'own';

/** The way an image job left one image (Go's hostapi.ImageWay). */
export type ImageWay = { image: string; way: ImageWayName };

/**
 * One image job flai serve runs, queued, or ran for the project (Go's hostapi.ImageJob): its id,
 * its images' refs, its state, the way it left each image once ended, its current step, when it
 * was asked for, started, and ended (RFC 3339), how it ended and why it failed, and its log's size.
 */
export type ImageJob = {
	/** Such as 20261113T093000Z-001. */
	id: string;
	images: string[];
	state: 'queued' | 'running' | 'ended';
	ways: ImageWay[];
	/** While it runs, its log's last "agent image:" line; once it has failed, the step that failed. */
	step?: string;
	requested: string;
	started?: string;
	ended?: string;
	outcome?: 'passed' | 'failed' | 'stopped';
	/** Why it failed or was stopped. */
	cause?: string;
	duration_ms?: number;
	/** 0 before it starts. */
	log_size: number;
};

/** What flai's images.get answers (Go's hostapi.ImagesState): GET /api/images. */
export type ImagesState = {
	/** The flai version the images must hold; empty when no image on the host was judged. */
	flai: string;
	images: AgentImage[];
	/** Why the images' state could not be told, such as docker not on flai serve's PATH. */
	unjudged?: string;
	running: ImageJob | null;
	/** The jobs that wait their turn, in order. */
	queued: ImageJob[];
	/** The last ended job of each set of images, newest first. */
	ended: ImageJob[];
};

/**
 * A read of an image job's log from an offset (Go's hostapi.ImageLog), as KlusterLog is of a
 * step's: at most 64 KB of text from offset, and next, the offset to ask from next, so a page
 * follows a running job live: GET /api/images?job=&offset=.
 */
export type ImageLog = {
	id: string;
	state: ImageJob['state'];
	/** Absent until the job has ended. */
	outcome?: ImageJob['outcome'];
	/** False while the job has no log yet, as when it is queued. */
	exists: boolean;
	/** Where text starts in the log. */
	offset: number;
	/** The offset to ask from next: offset and the bytes of text. */
	next: number;
	/** The log's size when it was read. */
	size: number;
	text: string;
	/** The log was shorter than the offset asked; text starts at 0. */
	reset?: boolean;
};

/** The job images.ensure queued (Go's hostapi.ImageStarted), and how many jobs wait or run before it. */
export type ImageStarted = ImageJob & { ahead: number };

/** What POST /api/images/retry answers: the job it started, then the images and jobs with it in them. */
export type ImagesRetried = ImagesState & { started: ImageStarted };

/**
 * The body of POST /api/images/retry's refusal while a job that makes the same images waits or
 * runs: flai's sentence, which names the job, and the job.
 */
export type ImagesBusy = { error: string; job: ImageJob };
