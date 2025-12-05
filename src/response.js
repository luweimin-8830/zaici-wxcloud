export function ok(data = null, message = 'ok') {
    return { code: 0, message, data };
}
export function fail(code = 500, message = 'error', data = null) {
    return { code, message, data };
}