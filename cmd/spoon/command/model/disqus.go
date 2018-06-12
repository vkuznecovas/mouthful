// package model contains all the required models for commands
package model

import "encoding/xml"

type Cauthor struct {
	XMLName      xml.Name      `xml:"author,omitempty" json:"author,omitempty"`
	Cemail       *Cemail       `xml:"http://disqus.com email,omitempty" json:"email,omitempty"`
	CisAnonymous *CisAnonymous `xml:"http://disqus.com isAnonymous,omitempty" json:"isAnonymous,omitempty"`
	Cname        *Cname        `xml:"http://disqus.com name,omitempty" json:"name,omitempty"`
	Cusername    *Cusername    `xml:"http://disqus.com username,omitempty" json:"username,omitempty"`
}

type Ccategory struct {
	XMLName        xml.Name    `xml:"category,omitempty" json:"category,omitempty"`
	AttrDsqSpaceid string      `xml:"http://disqus.com/disqus-internals id,attr"  json:",omitempty"`
	Cforum         *Cforum     `xml:"http://disqus.com forum,omitempty" json:"forum,omitempty"`
	CisDefault     *CisDefault `xml:"http://disqus.com isDefault,omitempty" json:"isDefault,omitempty"`
	Ctitle         *Ctitle     `xml:"http://disqus.com title,omitempty" json:"title,omitempty"`
}

type CcreatedAt struct {
	XMLName xml.Name `xml:"createdAt,omitempty" json:"createdAt,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type Cdisqus struct {
	XMLName                    xml.Name   `xml:"disqus,omitempty" json:"disqus,omitempty"`
	AttrXmlnsdsq               string     `xml:"xmlns dsq,attr"  json:",omitempty"`
	AttrXsiSpaceschemaLocation string     `xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"  json:",omitempty"`
	Attrxmlns                  string     `xml:"xmlns,attr"  json:",omitempty"`
	AttrXmlnsxsi               string     `xml:"xmlns xsi,attr"  json:",omitempty"`
	Ccategory                  *Ccategory `xml:"http://disqus.com category,omitempty" json:"category,omitempty"`
	Cpost                      []*Cpost   `xml:"http://disqus.com post,omitempty" json:"post,omitempty"`
	Cthread                    []*Cthread `xml:"http://disqus.com thread,omitempty" json:"thread,omitempty"`
}

type Cemail struct {
	XMLName xml.Name `xml:"email,omitempty" json:"email,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type Cforum struct {
	XMLName xml.Name `xml:"forum,omitempty" json:"forum,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type Cid struct {
	XMLName xml.Name `xml:"id,omitempty" json:"id,omitempty"`
}

type CipAddress struct {
	XMLName xml.Name `xml:"ipAddress,omitempty" json:"ipAddress,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type CisAnonymous struct {
	XMLName xml.Name `xml:"isAnonymous,omitempty" json:"isAnonymous,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type CisClosed struct {
	XMLName xml.Name `xml:"isClosed,omitempty" json:"isClosed,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type CisDefault struct {
	XMLName xml.Name `xml:"isDefault,omitempty" json:"isDefault,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type CisDeleted struct {
	XMLName xml.Name `xml:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type CisSpam struct {
	XMLName xml.Name `xml:"isSpam,omitempty" json:"isSpam,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type Clink struct {
	XMLName xml.Name `xml:"link,omitempty" json:"link,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type Cmessage struct {
	XMLName xml.Name `xml:"message,omitempty" json:"message,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type Cname struct {
	XMLName xml.Name `xml:"name,omitempty" json:"name,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type Cparent struct {
	XMLName        xml.Name `xml:"parent,omitempty" json:"parent,omitempty"`
	AttrDsqSpaceid string   `xml:"http://disqus.com/disqus-internals id,attr"  json:",omitempty"`
}

type Cpost struct {
	XMLName        xml.Name    `xml:"post,omitempty" json:"post,omitempty"`
	AttrDsqSpaceid string      `xml:"http://disqus.com/disqus-internals id,attr"  json:",omitempty"`
	Cauthor        *Cauthor    `xml:"http://disqus.com author,omitempty" json:"author,omitempty"`
	CcreatedAt     *CcreatedAt `xml:"http://disqus.com createdAt,omitempty" json:"createdAt,omitempty"`
	Cid            *Cid        `xml:"http://disqus.com id,omitempty" json:"id,omitempty"`
	CipAddress     *CipAddress `xml:"http://disqus.com ipAddress,omitempty" json:"ipAddress,omitempty"`
	CisDeleted     *CisDeleted `xml:"http://disqus.com isDeleted,omitempty" json:"isDeleted,omitempty"`
	CisSpam        *CisSpam    `xml:"http://disqus.com isSpam,omitempty" json:"isSpam,omitempty"`
	Cmessage       *Cmessage   `xml:"http://disqus.com message,omitempty" json:"message,omitempty"`
	Cparent        *Cparent    `xml:"http://disqus.com parent,omitempty" json:"parent,omitempty"`
	Cthread        []*Cthread  `xml:"http://disqus.com thread,omitempty" json:"thread,omitempty"`
}

type Cthread struct {
	XMLName        xml.Name    `xml:"thread,omitempty" json:"thread,omitempty"`
	AttrDsqSpaceid string      `xml:"http://disqus.com/disqus-internals id,attr"  json:",omitempty"`
	Cauthor        *Cauthor    `xml:"http://disqus.com author,omitempty" json:"author,omitempty"`
	Ccategory      *Ccategory  `xml:"http://disqus.com category,omitempty" json:"category,omitempty"`
	CcreatedAt     *CcreatedAt `xml:"http://disqus.com createdAt,omitempty" json:"createdAt,omitempty"`
	Cforum         *Cforum     `xml:"http://disqus.com forum,omitempty" json:"forum,omitempty"`
	Cid            *Cid        `xml:"http://disqus.com id,omitempty" json:"id,omitempty"`
	CipAddress     *CipAddress `xml:"http://disqus.com ipAddress,omitempty" json:"ipAddress,omitempty"`
	CisClosed      *CisClosed  `xml:"http://disqus.com isClosed,omitempty" json:"isClosed,omitempty"`
	CisDeleted     *CisDeleted `xml:"http://disqus.com isDeleted,omitempty" json:"isDeleted,omitempty"`
	Clink          *Clink      `xml:"http://disqus.com link,omitempty" json:"link,omitempty"`
	Cmessage       *Cmessage   `xml:"http://disqus.com message,omitempty" json:"message,omitempty"`
	Ctitle         *Ctitle     `xml:"http://disqus.com title,omitempty" json:"title,omitempty"`
}

type Ctitle struct {
	XMLName xml.Name `xml:"title,omitempty" json:"title,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}

type Cusername struct {
	XMLName xml.Name `xml:"username,omitempty" json:"username,omitempty"`
	SValue  string   `xml:",chardata" json:",omitempty"`
}
