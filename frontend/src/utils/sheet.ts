import { Principal, Sheet, SheetPayload, SheetSource } from "@/types";
import { hasProjectPermission, hasWorkspacePermission } from "../utils";

export const isSheetReadable = (sheet: Sheet, currentUser: Principal) => {
  // readable to
  // PRIVATE: the creator only
  // PROJECT: the creator and members in the project, workspace Owner and DBA
  // PUBLIC: everyone

  if (sheet.creator.id === currentUser.id) {
    // Always readable to the creator
    return true;
  }
  const { visibility } = sheet;
  if (visibility === "PRIVATE") {
    return false;
  }
  if (visibility === "PROJECT") {
    if (
      hasWorkspacePermission(
        "bb.permission.workspace.manage-project",
        currentUser.role
      )
    ) {
      return true;
    }

    const projectMemberList = sheet.project.memberList;
    return (
      projectMemberList.findIndex(
        (member) => member.principal.id === currentUser.id
      ) >= 0
    );
  }
  // visibility === "PUBLIC"
  return true;
};

export const isSheetWritable = (sheet: Sheet, currentUser: Principal) => {
  // writable to
  // PRIVATE: the creator only
  // PROJECT: the creator or project role can manage sheet, workspace Owner and DBA
  // PUBLIC: the creator only

  if (sheet.creator.id === currentUser.id) {
    // Always writable to the creator
    return true;
  }
  const { visibility } = sheet;
  if (visibility === "PRIVATE") {
    return false;
  }
  if (visibility === "PROJECT") {
    if (
      hasWorkspacePermission(
        "bb.permission.workspace.manage-project",
        currentUser.role
      )
    ) {
      return true;
    }

    const isCurrentUserProjectOwner = () => {
      const projectMemberList = sheet.project.memberList || [];
      const memberInProject = projectMemberList.find((member) => {
        return member.principal.id === currentUser.id;
      });

      return (
        memberInProject &&
        hasProjectPermission(
          "bb.permission.project.manage-sheet",
          memberInProject.role
        )
      );
    };
    return isCurrentUserProjectOwner();
  }
  // visibility === "PUBLIC"
  return false;
};

// getDefaultSheetPayloadWithSource gets the default payload with sheet source.
export const getDefaultSheetPayloadWithSource = (
  sheetSource: SheetSource
): SheetPayload => {
  if (sheetSource === "BYTEBASE") {
    // As we don't save any data for sheet from UI, return an empty payload.
    return {};
  }

  // Shouldn't reach this line.
  // For those sheet from VCS, we create and patch them in backend.
  return {};
};
