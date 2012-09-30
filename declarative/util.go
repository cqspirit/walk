// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

func InitWidget(d Widget, w walk.Widget, customInit func() error) error {
	var succeeded bool
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	// Widget
	if p := w.Parent(); p != nil {
		name, stretchFactor, row, rowSpan, column, columnSpan, contextMenu := d.CommonInfo()

		w.SetName(name)

		switch l := p.Layout().(type) {
		case *walk.BoxLayout:
			if stretchFactor < 1 {
				stretchFactor = 1
			}
			if err := l.SetStretchFactor(w, stretchFactor); err != nil {
				return err
			}

		case *walk.GridLayout:
			cs := columnSpan
			if cs < 1 {
				cs = 1
			}
			rs := rowSpan
			if rs < 1 {
				rs = 1
			}
			r := walk.Rectangle{column, row, cs, rs}

			if err := l.SetRange(w, r); err != nil {
				return err
			}
		}

		if contextMenu != nil {
			cm, err := walk.NewMenu()
			if err != nil {
				return err
			}
			if err := contextMenu.init(cm); err != nil {
				return err
			}
			w.SetContextMenu(cm)
		}
	}

	// Fonter
	if fonter, ok := d.(Fonter); ok {
		if f, err := fonter.Font_().Create(); err != nil {
			return err
		} else if f != nil {
			w.SetFont(f)
		}
	}

	// Container
	if dc, ok := d.(Container); ok {
		if wc, ok := w.(walk.Container); ok {
			layout, children := dc.ContainerInfo()

			if layout != nil {
				l, err := layout.Create()
				if err != nil {
					return err
				}

				if err := wc.SetLayout(l); err != nil {
					return err
				}
			}

			for _, child := range children {
				if err := child.Create(wc); err != nil {
					return err
				}
			}
		}
	}

	// Custom
	if customInit != nil {
		if err := customInit(); err != nil {
			return err
		}
	}

	succeeded = true

	return nil
}