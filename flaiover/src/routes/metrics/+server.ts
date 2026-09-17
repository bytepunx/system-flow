import { registry } from '$lib/server/metrics';

/** Prometheus exposition of the golden signals and build_info. */
export const GET = async () =>
	new Response(await registry.metrics(), { headers: { 'content-type': registry.contentType } });
