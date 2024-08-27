<template>
    <div>
        <BaseHeader />
        <div class="flex main-container">
            <BaseSide />
            <div w="full" py="4">
                <router-view />
            </div>
        </div>
        <Upgrade />
        <Notice />
        <Ad />
    </div>
</template>
<style>
.main-container {
    height: calc(100vh - var(--el-menu-item-height) - 5px);
    width: 100vw;
    background-color: var(--el-bg-color);
    overflow: hidden;
}
</style>