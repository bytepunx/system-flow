<script lang="ts">
	// Not destructured: the props are the type's members.
	type Props = { title: string; subtitle?: string };

	let props: Props = $props();
</script>

<h1>{props.title}</h1>
