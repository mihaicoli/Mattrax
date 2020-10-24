export interface LoginRequest {
  upn: string
  password: string
}

interface UserInformation {
  name?: string
  upn?: string
  org?: string
}

interface State {
  authToken: string
  user: UserInformation
}

export const state = (): State => ({
  authToken: sessionStorage.getItem('authToken') || '',
  user: {},
})

export const mutations = {
  setAuthToken(state: State, authToken: string) {
    state.authToken = authToken
  },

  setUserInformation(state: State, user: any) {
    state.user = user
  },
}

export const actions = {
  isAuthenticated(context: any): boolean {
    return context.state.authToken !== ''
  },

  populateUserInfomation(context: any) {
    if (context.state.authToken === '') {
      return
    }

    try {
      const base64Url = context.state.authToken.split('.')[1]
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
          })
          .join('')
      )
      const claims = JSON.parse(jsonPayload)

      const userInfo: UserInformation = {
        name: claims.name,
        upn: claims.sub,
        org: claims.org,
      }

      context.commit('setUserInformation', userInfo)
    } catch {}
  },

  login(context: any, user: LoginRequest) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/login', {
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

          const data = await res.json()
          sessionStorage.setItem('authToken', data.token)
          context.commit('setAuthToken', data.token)
          resolve()
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
