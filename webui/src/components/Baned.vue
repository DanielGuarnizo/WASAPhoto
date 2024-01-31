<script>
export default {
    props: ['msg']
}
</script>

<template>
    <div class="alert alert-warning" role="alert">
        <strong>Warning:</strong> You have been banned by <strong>{{ msg }}</strong> .
    </div>
</template>

<style>
</style>