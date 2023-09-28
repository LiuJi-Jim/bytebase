import Emittery from "emittery";
import { InjectionKey, Ref, computed, inject, provide, ref } from "vue";
import { useRoute } from "vue-router";
import { useChangelistStore, useProjectV1Store } from "@/store";
import { ComposedProject, unknownChangelist } from "@/types";
import {
  Changelist,
  Changelist_Change as Change,
} from "@/types/proto/v1/changelist_service";

export type ChangelistDetailEvents = Emittery<{
  "reorder-cancel": undefined;
  "reorder-confirm": undefined;
}>;

export type ChangelistDetailContext = {
  changelist: Ref<Changelist>;
  project: Ref<ComposedProject>;
  allowEdit: Ref<boolean>;
  reorderMode: Ref<boolean>;
  selectedChanges: Ref<Change[]>;
  showAddChangePanel: Ref<boolean>;
  isUpdating: Ref<boolean>;

  events: ChangelistDetailEvents;
};

export const KEY = Symbol(
  "bb.changelist.detail"
) as InjectionKey<ChangelistDetailContext>;

export const useChangelistDetailContext = () => {
  return inject(KEY)!;
};

export const provideChangelistDetailContext = () => {
  const route = useRoute();
  const name = computed(() => {
    return `projects/${route.params.projectName}/changelists/${route.params.changelistName}`;
  });
  const changelist = computed(() => {
    return (
      useChangelistStore().getChangelistByName(name.value) ??
      unknownChangelist()
    );
  });
  const project = computed(() => {
    return useProjectV1Store().getProjectByName(
      `projects/${route.params.projectName}`
    );
  });
  const allowEdit = computed(() => {
    return true; // TODO
  });

  const context: ChangelistDetailContext = {
    changelist,
    project,
    allowEdit,
    reorderMode: ref(false),
    selectedChanges: ref([]),
    showAddChangePanel: ref(false),
    isUpdating: ref(false),

    events: new Emittery(),
  };

  provide(KEY, context);

  return context;
};
