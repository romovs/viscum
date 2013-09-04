// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

//=====================================================================================================================
// Helper functions.
//
//=====================================================================================================================

package utils

import (
	"fmt"
	"reflect"
	"strings"
	log "github.com/cihub/seelog"
)

const FIELD_NAME_LENGTH int = 20


func StructPrint(s interface{}) string {

    val := reflect.ValueOf(s).Elem()
    
    info := ""; 
    
    for i := 0; i < val.NumField(); i++ {
        valueField := val.Field(i)
        typeField := val.Type().Field(i)

		field := fmt.Sprintf("%v", typeField.Name)
		
		if len(field) < FIELD_NAME_LENGTH {
			field = field + strings.Repeat(" ", FIELD_NAME_LENGTH - len(field))
		}
		
        info += fmt.Sprintf(field + ": %v\n", valueField.Interface())
    }
    
    return info
}


func LoadLogConfig(name string) {
	config := `<seelog type="adaptive" mininterval="2000000" maxinterval="10000000" critmsgcount="5">
				<outputs formatid="msg">
				<file path="` + name + `.log"/>
				</outputs>
				<formats>
				<format id="msg" format="%Time: %Msg%n"/>
				</formats>
				</seelog>`

	logger, _ := log.LoggerFromConfigAsBytes([]byte(config))

	log.ReplaceLogger(logger)
}
