import { ElMessage, ElMessageBox } from 'element-plus'

export function notifyError(message: string) {
    ElMessage({
        type: 'error',
        message
    })
}
export function notifySuccess(message: string) {
    ElMessage({
        type: 'success',
        message
    })
}
export function notifyInfo(message: string) {
    ElMessage({
        type: 'info',
        message
    })
}
export function buyVip() {
    return ElMessageBox.confirm('是否购买Vip版本？', '只有vip才能使用', {
        confirmButtonText: '立即购买',
        cancelButtonText: '再想一想',
        type: 'warning',
    })
}
export async function confirm(message: string, title: string, type: 'info' | 'warning' | 'error' = 'info'): Promise<number> {
    try {
        await ElMessageBox.confirm(message, title, {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: type,
        });
        return 1; // 用户点击了确认按钮
    } catch (error) {
        return 0; // 用户点击了取消或关闭对话框
    }
}
