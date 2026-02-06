<script>
  export let type = 'button';
  export let color = ''; // primary, link, info, success, warning, danger
  export let size = ''; // small, normal, medium, large
  export let loading = false;
  export let disabled = false;
  export let outlined = false;
  export let inverted = false;
  export let rounded = false;
  export let fullwidth = false;

  $: classes = [
    'button',
    color ? `is-${color}` : '',
    size ? `is-${size}` : '',
    loading ? 'is-loading' : '',
    outlined ? 'is-outlined' : '',
    inverted ? 'is-inverted' : '',
    rounded ? 'is-rounded' : '',
    fullwidth ? 'is-fullwidth' : '',
  ].filter(Boolean).join(' ');
</script>

<button
  {type}
  class={classes}
  disabled={disabled || loading}
  on:click
>
  <slot></slot>
</button>
