function StandardFetch(path: string, init?: RequestInit): Promise<Response> {
    return getStandardFetch(3, path, init);
}

function getStandardFetch(retries: number, path: string, init?: RequestInit): Promise<Response> {
    if (!init) {
        init = {};
    }

    // Always add credentials.
    init.credentials = "include";

    return fetch((process.env.REACT_APP_BACKEND_DOMAIN || "") + "/" + path, init)
    .then(response => {
        if (response.status === 401 || response.status === 403) {
            // This is probably not a good pattern, right?
            window.location.href = "/sign-in?state=unauthorized";
            return Promise.reject(); 
        } else if (response.status === 504) {
            // To help with Aurora waking up.
            if (retries > 0) {
                return getStandardFetch(--retries, path, init);
            }
        }
        return response;
    });

}

export { StandardFetch };