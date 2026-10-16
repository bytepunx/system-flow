<script lang="ts">
	// Typed by a local alias given to $props as a type argument.
	type Variant = 'primary' | 'secondary';
	type Props = { variant?: Variant; disabled?: boolean };

	let { variant = 'primary', disabled } = $props<Props>();
</script>

<button class={variant} {disabled}>ok</button>
