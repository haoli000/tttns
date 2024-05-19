/*
Copyright © 2024 Hao Li <mr.hao.li@gmail.com>

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

// For an introduction to this cli tool, see the [docs].
// To decode the dumped CDR contents, another cli tool [xchf] can be used.
//
// [docs]: https://github.com/haoli000/tttns/tree/main/docs
// [xchf]: https://github.com/haoli000/xchf
package main

import "github.com/haoli000/tttns/cmd"

func main() {
	cmd.Execute()
}
