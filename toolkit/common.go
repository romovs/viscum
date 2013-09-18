// Copyright (C) 2013 Roman Ovseitsev <romovs@gmail.com>
// This software is distributed under GNU GPL v2. See LICENSE file.

package toolkit

import (

)

// button click handler. parameter is ignored for non BS_TOGGLE buttons.
type clickHandler func(bool)