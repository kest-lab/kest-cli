export const authConfig = {
    routes: {
        login: '/login',
        register: '/register',
        afterLogin: '/',
        afterLogout: '/login',
    },
    cookieKeys: {
        accessToken: 'kest_access_token',
        refreshToken: 'kest_refresh_token',
    },
};

export const mockUsers: any[] = [];
