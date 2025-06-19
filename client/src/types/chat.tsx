export type ChatMessage = {
    type: string
    payload: {
        user: {
            id: string
            name: string,
        }
        message: string
    }
}