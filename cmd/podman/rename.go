package main

import (
    "github.com/containers/buildah"
    lu "github.com/containers/libpod/pkg/util"
    "github.com/containers/storage"
    "github.com/pkg/errors"
    "github.com/containers/libpod/cmd/podman/cliconfig"
    "github.com/spf13/cobra"
)

var (
    renameCommand     cliconfig.RenameValues
    renameDescription = "Rename a container."

    _renameCommand = &cobra.Command{
        Use:   "rename",
        Short: "Rename one or more containers.",
        Long:  renameDescription,
        RunE: func(cmd *cobra.Command, args []string) error {
            renameCommand.InputArgs = args
            renameCommand.GlobalFlags = MainGlobalOpts
            return renameCmd(&renameCommand)
        },
        Example: "CONTAINER NEW_NAME",
    }
)

func init() {
        renameCommand.Command = _renameCommand
        rootCmd.AddCommand(renameCommand.Command)
}

// Define the entrypoint command
func renameCmd(c *cliconfig.RenameValues) error {

    var builder *buildah.Builder

    args := c.InputArgs
    name := args[0]
    newName := args[1]

    store, err := storage.GetStore(options)(c)
    if err != nil {
            return err
    }

    builder, err = openBuilder(getContext(), store, name)
    if err != nil {
            return errors.Wrapf(err, "error reading build container %q", name)
    }

    oldName := builder.Container
    if oldName == newName {
            return errors.Errorf("renaming a container with the same name as its current name")
    }

    if build, err := openBuilder(getContext(), store, newName); err == nil {
            return errors.Errorf("The container name %q is already in use by container %q", newName, build.ContainerID)
    }

    err = store.SetNames(builder.ContainerID, []string{newName})
    if err != nil {
            return errors.Wrapf(err, "error renaming container %q to the name %q", oldName, newName)
    }
    builder.Container = newName
    return builder.Save()
}
