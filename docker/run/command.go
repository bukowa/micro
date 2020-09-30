/*
Copyright © 2020 Mateusz Kurowski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package run

import (
	"os/exec"
)

type Command struct {
	Image string
	Flags Flags
	Args  []string
}

func NewCommand(image string, args []string, flags ...Flag) *Command {
	if args == nil {
		args = []string{}
	}
	return &Command{
		Image: image,
		Flags: NewFlags(flags...),
		Args:  args,
	}
}

func (c *Command) Arg(a ...string) {
	c.Args = append(c.Args, a...)
}

func (c *Command) Command() *exec.Cmd {
	cmd := exec.Command("docker")
	cmd.Args = append(append(append(append(append(cmd.Args, "run")), c.Flags...), c.Image), c.Args...)
	return cmd
}
