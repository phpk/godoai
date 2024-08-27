export async function OpenDirDialog(){
    if((window as any).go) {
        //(window as any).go.OpenDirDialog();
        return (window as any)['go']['main']['App']['OpenDirDialog']();
    }else {
        return ""
    }
}
export async function RestartApp(){
    if((window as any).go) {
        return (window as any)['go']['main']['App']['RestartApp']();
    }else {
        window.location.reload();
    }
}