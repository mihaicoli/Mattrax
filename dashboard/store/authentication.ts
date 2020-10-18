export interface LoginRequest {
  upn: string
  password: string
}

interface State {
  authToken: string
  user: {
    name?: string
    upn?: string
    org?: string
  }
}

export const state = (): State => ({
  authToken: '1', // TODO: TEMP Value
  user: {},
})

export const mutations = {
  setAuthToken(state: State, authToken: string) {
    state.authToken = authToken
  },

  setTEMPUser(state: State, user: any) {
    state.user = user
  },
}

export const actions = {
  isAuthenticated(context: any): boolean {
    return context.state.authToken !== ''
  },

  login(context: any, user: LoginRequest) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/010341b0-913f-42a7-86eb-23f57bb9a0fe', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(user),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('The login request was rejected'))
            return
          }

          // TODO: Replace with getting user information from JWT authToken not response
          const data = await res.json()
          console.log(data)
          context.commit('setTEMPUser', data)

          context.commit('setAuthToken', 'test')
          resolve()
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
