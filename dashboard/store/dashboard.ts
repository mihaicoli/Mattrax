interface State {
  error: Error | null
}

export const state = (): State => ({
  error: null,
})

export const mutations = {
  setError(state: State, error: Error) {
    state.error = error
  },
  clearError(state: State) {
    state.error = null
  },
}
