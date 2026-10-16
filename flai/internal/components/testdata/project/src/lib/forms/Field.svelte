<script lang="ts">
	// Typed by a local interface that extends another; rest takes the others.
	interface Base {
		id: string;
	}

	interface Props extends Base {
		value?: string;
		level: 1 | 2 | 3;
		validate(value: string): boolean;
	}

	let { value = '', level, ...rest }: Props = $props();
</script>

<input {value} data-level={level} {...rest} />
