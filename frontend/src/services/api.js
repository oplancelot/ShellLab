export const GetQuestDetail = (entry) => {
    console.log(`[API] Fetching Quest Detail for: ${entry}`);
    if (window?.go?.main?.App?.GetQuestDetail) {
        return window.go.main.App.GetQuestDetail(entry)
            .then(res => {
                console.log(`[API] Received Quest Detail for ${entry}:`, res);
                return res;
            })
            .catch(err => {
                console.error(`[API] Failed to get Quest Detail for ${entry}:`, err);
                throw err;
            });
    }
    console.warn(`[API] GetQuestDetail not found in Wails App!`);
    return Promise.resolve(null)
}

export const GetCreatureDetail = (entry) => {
    console.log(`[API] Fetching Creature Detail for: ${entry}`);
    if (window?.go?.main?.App?.GetCreatureDetail) {
        return window.go.main.App.GetCreatureDetail(entry);
    }
    return Promise.resolve(null)
}

export const GetItemDetail = (entry) => {
    console.log(`[API] Fetching Item Detail for: ${entry}`);
    if (window?.go?.main?.App?.GetItemDetail) {
        return window.go.main.App.GetItemDetail(entry);
    }
    return Promise.resolve(null)
}
