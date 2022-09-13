function StandardFetch(path: string, init?: RequestInit): Promise<Response> {
    if (!init) {
        init = {}
    }

    // Always add credentials.
    init.credentials = "include"

    return fetch((process.env.REACT_APP_BACKEND_DOMAIN || "") + "/" + path, init)
    .then(response => {
        if (response.status === 401 || response.status === 403) {
            const next = encodeURIComponent(window.location.pathname + window.location.search)
            // Is setting the location directly a good pattern?
            window.location.href = "/sign-in?state=" + response.status + "&next=" + next
            return Promise.reject()
        }
        return response
    })

}

export { StandardFetch }