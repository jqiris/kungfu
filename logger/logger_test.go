/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package logger

import (
	"testing"
)

func TestLoggerColor(t *testing.T) {
	//fmt.Printf("\033[1;37;41m%s\033[0m\n", "Red.")
	//d := color.New(color.FgHiYellow)
	//_, err := d.Printf("hello world")
	//if err != nil {
	//	t.Fatal(err)
	//}

	Info("hello world")
	n := defLogger.WithSuffix("welcome")
	n.Info("hello world2")
	select {}

}
