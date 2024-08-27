<template>
    <div class="flex main-container">
        <router-view />
    </div>
</template>

<style scoped>
.main-container {
    height: 100vh;
    width: 100vw;
    background-color: var(--el-bg-color);
    overflow: hidden;
}
</style>