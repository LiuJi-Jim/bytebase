import axios from "axios";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { authServiceClient } from "@/grpcweb";
import { unknownUser } from "@/types";
import { userBindingPrefix } from "@/types";
import type {
  LoginRequest,
  LoginResponse,
  User,
} from "@/types/proto/v1/auth_service";
import { UserType } from "@/types/proto/v1/auth_service";
import { getIntCookie } from "@/utils";
import { useUserStore, useWorkspaceV1Store } from ".";

export const useAuthStore = defineStore("auth_v1", () => {
  const userStore = useUserStore();
  const workspaceStore = useWorkspaceV1Store();
  const currentUserId = ref<number | undefined>();

  const currentUser = computed(() => {
    if (currentUserId.value) {
      return userStore.getUserById(`${currentUserId.value}`) ?? unknownUser();
    }
    return unknownUser();
  });

  const currentRolesInWorkspace = computed(() => {
    return workspaceStore.findRolesByMember({
      member: `${userBindingPrefix}${currentUser.value.email}`,
      ignoreGroup: false,
    });
  });

  const isLoggedIn = () => {
    return getUserIdFromCookie() != undefined;
  };

  const getUserIdFromCookie = () => {
    return getIntCookie("user");
  };

  const login = async (request: Partial<LoginRequest>) => {
    const {
      data: { mfaTempToken },
    } = await axios.post<LoginResponse>("/v1/auth/login", request);
    if (mfaTempToken) {
      return mfaTempToken;
    }

    await restoreUser();
  };

  const signup = async (request: Partial<User>) => {
    await authServiceClient.createUser({
      user: {
        email: request.email,
        title: request.name,
        password: request.password,
        userType: UserType.USER,
      },
    });
    await login({
      email: request.email,
      password: request.password,
      web: true,
    });
    await workspaceStore.fetchIamPolicy();
  };

  const logout = async () => {
    try {
      await axios.post("/v1/auth/logout");
      currentUserId.value = undefined;
    } catch {
      // nothing
    }
  };

  const restoreUser = async () => {
    currentUserId.value = getUserIdFromCookie();
    if (currentUserId.value) {
      await useUserStore().getOrFetchUserById(
        String(currentUserId.value),
        true // silent
      );
      await workspaceStore.fetchIamPolicy();
    }
  };

  const refreshUserIfNeeded = async (name: string) => {
    if (name === currentUser.value.name) {
      await useUserStore().fetchUser(
        name,
        true // silent
      );
    }
  };

  return {
    currentUser,
    currentUserId,
    currentRolesInWorkspace,
    isLoggedIn,
    getUserIdFromCookie,
    login,
    signup,
    logout,
    restoreUser,
    refreshUserIfNeeded,
  };
});

export const useCurrentRoles = () => {
  const authStore = useAuthStore();
  return computed(() => authStore.currentRolesInWorkspace);
};

export const useCurrentUserV1 = () => {
  const authStore = useAuthStore();
  return computed(() => authStore.currentUser);
};

export const useIsLoggedIn = () => {
  const store = useAuthStore();
  return computed(() => store.isLoggedIn() && store.currentUser.name !== "");
};
