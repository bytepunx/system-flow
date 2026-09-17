/** Short human age from seconds: 0m, 45m, 3h, 2d. */
export function age(seconds: number): string {
	if (seconds < 60) return '0m';
	if (seconds < 3600) return `${Math.floor(seconds / 60)}m`;
	if (seconds < 172800) return `${Math.floor(seconds / 3600)}h`;
	return `${Math.floor(seconds / 86400)}d`;
}
