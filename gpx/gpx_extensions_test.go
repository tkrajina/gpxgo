package gpx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadExtension(t *testing.T) {
	t.Parallel()

	original, reparsed, fromJSON := loadAndReparseFile(t, "../test_files/gpx1.1_with_extensions.gpx")

	byts, err := reparsed.ToXml(ToXmlParams{Indent: true})
	assert.Nil(t, err)
	fmt.Println(string(byts))

	/*
	   <extensions>
	       <ext:aaa ext:jjj="kkk">bbb</ext:aaa>hhh
	       <ext:ccc>
	           <ext:ddd ext:lll="mmm" ext:nnn="ooo">
	               <ext:fff>ggg</ext:fff>
	           </ext:ddd>
	       </ext:ccc>
	   </extensions>
	*/

	for n, g := range []GPX{*original, *reparsed, *fromJSON} {
		fmt.Printf("gpx #%d\n", n)

		exts := []Extension{
			g.MetadataExtensions,
			g.Routes[0].Points[0].Extensions,
			g.Waypoints[0].Extensions,
			g.Tracks[0].Segments[0].Points[0].Extensions,
		}

		for _, ext := range exts {
			assert.Equal(t, 2, len(ext))
			assert.Equal(t, "bbb", ext[0].Data)
			assert.Equal(t, 1, len(ext[0].Attrs), "%#v", ext[0].Attrs)
			assert.Equal(t, "kkk", ext[0].GetAttrOrEmpty("jjj"))
			assert.Equal(t, "aaa", ext[0].NameLocal)
			//assert.Equal(t, "gpx.py", ext.Nodes[0].SpaceName())
			assert.Equal(t, 1, len(ext[1].Nodes))
			assert.Equal(t, 0, len(ext[1].Attrs))
			assert.Equal(t, "mmm", ext[1].Nodes[0].GetAttrOrEmpty("lll"))
			assert.Equal(t, "ooo", ext[1].Nodes[0].GetAttrOrEmpty("nnn"))
			assert.Equal(t, "ggg", ext[1].Nodes[0].Nodes[0].Data)
		}
	}
}

func TestExtensionWithoutNamespace(t *testing.T) {
	t.Parallel()

	original, err := ParseString(`<gpx version="1.1" creator="nawagers" xmlns="http://www.topografix.com/GPX/1/1">
	<metadata>
		<extensions>
			<aaa jjj="kkk">bbb</aaa>hhh
			<ccc>
				<ddd lll="mmm" nnn="ooo">
					<fff>ggg</fff>
				</ddd>
			</ccc>
		</extensions>
	</metadata>
</gpx>`)
	assert.Nil(t, err)
	assert.NotNil(t, original)

	if t.Failed() {
		t.FailNow()
	}

	reparsed, err := reparse(*original)
	assert.Nil(t, err)

	for _, g := range []GPX{*original, *reparsed} {
		ext := g.MetadataExtensions
		assert.Equal(t, 2, len(ext))
		assert.Equal(t, "bbb", ext[0].Data)
		assert.Equal(t, 1, len(ext[0].Attrs), "%#v", ext[0].Attrs)
		assert.Equal(t, "kkk", ext[0].GetAttrOrEmpty("jjj"))
		assert.Equal(t, "aaa", ext[0].NameLocal)
		//assert.Equal(t, "gpx.py", ext.Nodes[0].SpaceName())
		assert.Equal(t, 1, len(ext[1].Nodes))
		assert.Equal(t, 0, len(ext[1].Attrs))
		assert.Equal(t, "mmm", ext[1].Nodes[0].GetAttrOrEmpty("lll"))
		assert.Equal(t, "ooo", ext[1].Nodes[0].GetAttrOrEmpty("nnn"))
		assert.Equal(t, "ggg", ext[1].Nodes[0].Nodes[0].Data)
	}
}

func TestNodesSubnodesAndAttrs(t *testing.T) {
	t.Parallel()

	var node ExtensionNode

	assert.Equal(t, 0, len(node.Attrs))
	node.SetAttr("xxx", "yyy")
	assert.Equal(t, 1, len(node.Attrs))
	{
		val, found := node.GetAttr("xxx")
		assert.True(t, found)
		assert.Equal(t, "yyy", val)
	}

	assert.Equal(t, 0, len(node.Nodes))
	node.GetOrCreateNode("aaa").Data = "aaa data"
	assert.Equal(t, 1, len(node.Nodes))
	assert.Equal(t, 0, len(node.Nodes[0].Attrs))

	assert.Equal(t, &node.Nodes[0], node.GetOrCreateNode("aaa"))

	// fmt.Println(string(node.debugXMLChunk()))
	node.GetOrCreateNode("aaa").SetAttr("aaa", "bbb")
	// fmt.Println(string(node.debugXMLChunk()))
	assert.Equal(t, 1, len(node.Nodes[0].Attrs))
	assert.Equal(t, "aaa", node.Nodes[0].Attrs[0].NameLocal)
	assert.Equal(t, "bbb", node.Nodes[0].Attrs[0].Value)

	// fmt.Println(string(node.debugXMLChunk()))
	node.GetOrCreateNode("aaa", "bbb").SetAttr("aaa", "bbb")
	// fmt.Println(string(node.debugXMLChunk()))
	assert.Equal(t, 1, len(node.Nodes))
	assert.Equal(t, 1, len(node.Nodes[0].Nodes))
	assert.Equal(t, "aaa", node.Nodes[0].Nodes[0].Attrs[0].NameLocal)
	assert.Equal(t, "bbb", node.Nodes[0].Nodes[0].Attrs[0].Value)
}

func TestExtensionNodesAndAttrs(t *testing.T) {
	t.Parallel()

	var ext Extension
	assert.Equal(t, 0, len(ext))
	ext.GetOrCreateNode(NoNamespace, "aaa").Data = "aaa data"
	assert.Equal(t, 1, len(ext))
	assert.Equal(t, 0, len(ext[0].Attrs))
	ext.GetOrCreateNode(NoNamespace, "aaa").SetAttr("aaa", "bbb")
	assert.Equal(t, 1, len(ext[0].Attrs))
	assert.Equal(t, "aaa", ext[0].Attrs[0].NameLocal)
	assert.Equal(t, "bbb", ext[0].Attrs[0].Value)

	// fmt.Println(string(ext.debugXMLChunk()))
	ext.GetOrCreateNode(NoNamespace, "aaa", "bbb").SetAttr("aaa", "bbb")
	// fmt.Println(string(ext.debugXMLChunk()))

	{
		// fmt.Println("a", string(ext.debugXMLChunk()))
		n1 := ext.GetOrCreateNode(NoNamespace, "aaa", "bbb")
		// fmt.Println("b", string(ext.debugXMLChunk()))
		n2 := &ext[0].Nodes[0]
		// fmt.Println("c", string(ext.debugXMLChunk()))
		assert.Equal(t, fmt.Sprintf("%p", n1), fmt.Sprintf("%p", n2))
	}

	assert.Equal(t, 1, len(ext))
	assert.Equal(t, 1, len(ext[0].Nodes))
	assert.Equal(t, "aaa", ext[0].Nodes[0].Attrs[0].NameLocal)
	assert.Equal(t, "bbb", ext[0].Nodes[0].Attrs[0].Value)
}

func TestCreateExtensionWithoutNamespace(t *testing.T) {
	t.Parallel()

	var original GPX
	// fmt.Println("1:", string(original.MetadataExtensions.debugXMLChunk()))
	original.MetadataExtensions.GetOrCreateNode(NoNamespace, "aaa", "bbb", "ccc").Data = "ccc data"
	// fmt.Println("2:", string(original.MetadataExtensions.debugXMLChunk()))
	assert.Equal(t, 1, len(original.MetadataExtensions))
	assert.Equal(t, "aaa", original.MetadataExtensions[0].NameLocal)
	assert.Equal(t, "bbb", original.MetadataExtensions[0].Nodes[0].NameLocal)
	assert.Equal(t, 0, len(original.MetadataExtensions[0].Nodes[0].Attrs), "attrs=%#v", original.MetadataExtensions[0].Nodes[0].Attrs)
	original.MetadataExtensions.GetOrCreateNode(NoNamespace, "aaa", "bbb").SetAttr("key", "value")
	// fmt.Println("3:", string(original.MetadataExtensions.debugXMLChunk()))
	assert.Equal(t, 1, len(original.MetadataExtensions[0].Nodes[0].Attrs), "attrs=%#v", original.MetadataExtensions[0].Nodes[0].Attrs)
	if t.Failed() {
		t.FailNow()
	}

	assert.Equal(t, "aaa", original.MetadataExtensions[0].NameLocal)
	assert.Equal(t, "bbb", original.MetadataExtensions[0].Nodes[0].NameLocal)
	assert.Equal(t, 1, len(original.MetadataExtensions[0].Nodes[0].Attrs), "attrs=%#v", original.MetadataExtensions[0].Nodes[0].Attrs)
	assert.Equal(t, "key", original.MetadataExtensions[0].Nodes[0].Attrs[0].NameLocal)
	assert.Equal(t, "value", original.MetadataExtensions[0].Nodes[0].Attrs[0].Value)

	val, found := original.MetadataExtensions.GetOrCreateNode(NoNamespace, "aaa", "bbb").GetAttr("key")
	assert.True(t, found)
	assert.Equal(t, "value", val)

	reparsed, err := reparse(original)
	assert.Nil(t, err)

	for _, g := range []GPX{original, *reparsed} {
		byts, err := g.ToXml(ToXmlParams{Indent: true})
		assert.Nil(t, err)
		expected := `<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="https://github.com/tkrajina/gpxgo" xmlns="http://www.topografix.com/GPX/1/1" >
       <metadata>
               <author></author>
               <extensions>
                       <aaa>
                               <bbb key="value">
                                       <ccc>ccc data</ccc>
                               </bbb>
                       </aaa>
               </extensions>
       </metadata>
</gpx>`
		assertLinesEquals(t, expected, string(byts))
	}
}

func TestCreateMetadataExtensionWithNamespace(t *testing.T) {
	t.Parallel()

	var original GPX
	original.RegisterNamespace("ext", "http://trla.baba.lan")
	original.MetadataExtensions.GetOrCreateNode("http://trla.baba.lan", "aaa", "bbb", "ccc").Data = "ccc data"

	assert.Equal(t, "http://trla.baba.lan", original.Attrs.NamespaceAttributes["xmlns"]["ext"].Value)
	assert.NotEmpty(t, original.Attrs.NamespaceAttributes["xmlns"]["ext"].replacement)

	original.MetadataExtensions.GetOrCreateNode("http://trla.baba.lan", "aaa", "bbb").SetAttr("key", "value")
	val, found := original.MetadataExtensions.GetOrCreateNode("http://trla.baba.lan", "aaa", "bbb").GetAttr("key")
	assert.True(t, found)
	assert.Equal(t, "value", val)

	reparsed, err := reparse(original)
	assert.Nil(t, err)

	rereparsed, err := reparse(*reparsed)
	assert.Nil(t, err)

	// fmt.Println(string(original.MetadataExtensions.debugXMLChunk()))
	// fmt.Println(string(reparsed.MetadataExtensions.debugXMLChunk()))
	// assert.Equal(t, original.MetadataExtensions.debugXMLChunk(), reparsed.MetadataExtensions.debugXMLChunk())
	assert.Equal(t, original.MetadataExtensions, reparsed.MetadataExtensions)

	assert.Equal(t, 1, len(original.Attrs.NamespaceAttributes))
	assert.Equal(t, len(original.Attrs.NamespaceAttributes), len(reparsed.Attrs.NamespaceAttributes))
	assert.Equal(t, original.Attrs.NamespaceAttributes["xmlns"]["ext"].Local, reparsed.Attrs.NamespaceAttributes["xmlns"]["ext"].Local)
	assert.Equal(t, original.Attrs.NamespaceAttributes["xmlns"]["ext"].Space, reparsed.Attrs.NamespaceAttributes["xmlns"]["ext"].Space)
	assert.Equal(t, original.Attrs.NamespaceAttributes["xmlns"]["ext"].Value, reparsed.Attrs.NamespaceAttributes["xmlns"]["ext"].Value)

	// assert.Equal(t, 1, len(reparsed.MetadataExtensions.Nodes))
	// assert.Equal(t, len(original.MetadataExtensions.Nodes), len(reparsed.MetadataExtensions.Nodes))
	// assert.Equal(t, original.MetadataExtensions.XMLName, reparsed.MetadataExtensions.XMLName)
	// assert.Equal(t, original.MetadataExtensions.Nodes[0], reparsed.MetadataExtensions.Nodes[0])
	// assert.Equal(t, original.MetadataExtensions.Attrs, reparsed.MetadataExtensions.Attrs)
	// assert.Equal(t, original.MetadataExtensions.Data, reparsed.MetadataExtensions.Data)
	assert.Equal(t, original.MetadataExtensions, reparsed.MetadataExtensions)

	if t.Failed() {
		t.FailNow()
	}

	for n, g := range []GPX{original, *reparsed, *rereparsed} {
		fmt.Printf("Test %d\n", n)

		node, found := g.MetadataExtensions.GetNode(AnyNamespace, "aaa")
		assert.True(t, found)
		assert.NotNil(t, node)

		node, found = g.MetadataExtensions.GetNode(NamespaceURL("http://trla.baba.lan"), "aaa")
		assert.True(t, found)
		assert.NotNil(t, node)
		assert.Equal(t, "http://trla.baba.lan", node.NameSpace)

		node, found = node.GetNode("bbb")
		assert.True(t, found)
		assert.NotNil(t, node)

		assert.Equal(t, "http://trla.baba.lan", node.NameSpace)

		byts, err := g.ToXml(ToXmlParams{Indent: true})
		assert.Nil(t, err)
		expected := `<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="https://github.com/tkrajina/gpxgo" xmlns="http://www.topografix.com/GPX/1/1" xmlns:ext="http://trla.baba.lan">
       <metadata>
               <author></author>
               <extensions>
                       <ext:aaa>
                               <ext:bbb ext:key="value">
                                       <ext:ccc>ccc data</ext:ccc>
                               </ext:bbb>
                       </ext:aaa>
               </extensions>
       </metadata>
</gpx>`
		assertLinesEquals(t, expected, string(byts))
	}
}

func TestCreateExtensionWithNamespace(t *testing.T) {
	t.Parallel()

	var original GPX
	original.RegisterNamespace("ext", "http://trla.baba.lan")
	original.Extensions.GetOrCreateNode("http://trla.baba.lan", "aaa", "bbb", "ccc").Data = "ccc data"

	assert.Equal(t, "http://trla.baba.lan", original.Attrs.NamespaceAttributes["xmlns"]["ext"].Value)
	assert.NotEmpty(t, original.Attrs.NamespaceAttributes["xmlns"]["ext"].replacement)

	original.Extensions.GetOrCreateNode("http://trla.baba.lan", "aaa", "bbb").SetAttr("key", "value")
	val, found := original.Extensions.GetOrCreateNode("http://trla.baba.lan", "aaa", "bbb").GetAttr("key")
	assert.True(t, found)
	assert.Equal(t, "value", val)

	reparsed, err := reparse(original)
	assert.Nil(t, err)

	rereparsed, err := reparse(*reparsed)
	assert.Nil(t, err)

	// fmt.Println(string(original.Extensions.debugXMLChunk()))
	// fmt.Println(string(reparsed.Extensions.debugXMLChunk()))
	// assert.Equal(t, original.Extensions.debugXMLChunk(), reparsed.Extensions.debugXMLChunk())
	assert.Equal(t, original.Extensions, reparsed.Extensions)

	assert.Equal(t, 1, len(original.Attrs.NamespaceAttributes))
	assert.Equal(t, len(original.Attrs.NamespaceAttributes), len(reparsed.Attrs.NamespaceAttributes))
	assert.Equal(t, original.Attrs.NamespaceAttributes["xmlns"]["ext"].Local, reparsed.Attrs.NamespaceAttributes["xmlns"]["ext"].Local)
	assert.Equal(t, original.Attrs.NamespaceAttributes["xmlns"]["ext"].Space, reparsed.Attrs.NamespaceAttributes["xmlns"]["ext"].Space)
	assert.Equal(t, original.Attrs.NamespaceAttributes["xmlns"]["ext"].Value, reparsed.Attrs.NamespaceAttributes["xmlns"]["ext"].Value)

	assert.Equal(t, 1, len(reparsed.Extensions))
	assert.Equal(t, len(original.Extensions), len(reparsed.Extensions))
	// assert.Equal(t, original.Extensions.XMLName, reparsed.Extensions.XMLName)
	assert.Equal(t, original.Extensions[0], reparsed.Extensions[0])
	// assert.Equal(t, original.Extensions.Attrs, reparsed.Extensions.Attrs)
	// assert.Equal(t, original.Extensions.Data, reparsed.Extensions.Data)
	assert.Equal(t, original.Extensions, reparsed.Extensions)

	if t.Failed() {
		t.FailNow()
	}

	for n, g := range []GPX{original, *reparsed, *rereparsed} {
		fmt.Printf("Test %d\n", n)

		node, found := g.Extensions.GetNode(AnyNamespace, "aaa")
		assert.True(t, found)
		assert.NotNil(t, node)

		node, found = g.Extensions.GetNode(NamespaceURL("http://trla.baba.lan"), "aaa")
		assert.True(t, found)
		assert.NotNil(t, node)
		assert.Equal(t, "http://trla.baba.lan", node.NameSpace)

		node, found = node.GetNode("bbb")
		assert.True(t, found)
		assert.NotNil(t, node)

		assert.Equal(t, "http://trla.baba.lan", node.NameSpace)

		byts, err := g.ToXml(ToXmlParams{Indent: true})
		assert.Nil(t, err)
		expected := `<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="https://github.com/tkrajina/gpxgo" xmlns="http://www.topografix.com/GPX/1/1" xmlns:ext="http://trla.baba.lan">
		<metadata>
  			<author></author>
        </metadata>
		<extensions>
				<ext:aaa>
						<ext:bbb ext:key="value">
								<ext:ccc>ccc data</ext:ccc>
						</ext:bbb>
				</ext:aaa>
		</extensions>
</gpx>`
		assertLinesEquals(t, expected, string(byts))
	}
}

func TestGarminExtensions(t *testing.T) {
	t.Parallel()

	original, reparsed, fromJSON := loadAndReparseFile(t, "../test_files/gpx_with_garmin_extension.gpx")
	if t.Failed() {
		t.FailNow()
	}

	for n, g := range []GPX{*original, *reparsed, *fromJSON} {
		fmt.Printf("Test %d\n", n)

		xml, err := g.ToXml(ToXmlParams{})
		assert.Nil(t, err)
		fmt.Println("xml:\n", string(xml))

		jsn, err := g.ToJSON(ToJSONParams{Indent: true})
		assert.Nil(t, err)
		fmt.Println("json:\n", string(jsn))

		assert.Nil(t, err)
		assert.Contains(t, string(xml), "<gpxtpx:TrackPointExtension>")
		assert.Contains(t, string(xml), "<gpxtpx:hr>171</gpxtpx:hr>")
	}
}
